package storage

// minio client wrapper

import (
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/maxkrivich/ProjectX/configs"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/policy"
)

type StorageClient struct {
	client  *minio.Client
	configs configs.MinioConfig
}

type StorageNotificationEvent minio.NotificationEvent

type HandlerFunction func(<-chan StorageNotificationEvent, <-chan StorageNotificationEvent, chan struct{}, *sync.WaitGroup)

func DialStorage() *StorageClient {
	var sc StorageClient
	var err error
	sc.configs = configs.GetConfig().MinioConfig
	sc.client, err = minio.New(sc.configs.Endpoint, sc.configs.AccessKeyID, sc.configs.SecretAccessKey, sc.configs.UseSSL)
	if err != nil {
		log.Fatal(err)
	}
	sc.setupBuckets()
	sc.setupBucketPolicy()
	return &sc
}

func (self *StorageClient) setupBuckets() {
	if err := self.client.MakeBucket(self.configs.FileBucketName, ""); err != nil {

		exists, err := self.client.BucketExists(self.configs.FileBucketName)
		if err == nil && exists {
			log.Printf("Bucket '%s' is already exists\n", self.configs.FileBucketName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("Successfully created '%s'\n", self.configs.FileBucketName)
}

func (self *StorageClient) setupBucketPolicy() {
	err := self.client.SetBucketPolicy(self.configs.FileBucketName, "", policy.BucketPolicyReadWrite)
	if err != nil {
		log.Fatal(err)
	}
}

func (self *StorageClient) PresignedPutObject(bucketName, objectName string, expire time.Duration) (u *url.URL, err error) {
	return self.client.PresignedPutObject(bucketName, objectName, expire)
}

func (self *StorageClient) PresignedGetObject(bucketName, objectName string, expires time.Duration, reqParams url.Values) (u *url.URL, err error) {
	return self.client.PresignedGetObject(bucketName, objectName, expires, reqParams)
}

func (self *StorageClient) RemoveObject(bucketName, objectName string) error {
	return self.client.RemoveObject(bucketName, objectName)
}

func (self *StorageClient) ListenBucketNotification(bucketName string, handlers ...HandlerFunction) {
	doneCh := make(chan struct{})
	objCreate, objRemove := make(chan StorageNotificationEvent), make(chan StorageNotificationEvent)
	defer close(doneCh)
	defer close(objRemove)
	defer close(objCreate)

	var wg sync.WaitGroup

	for _, f :=range handlers {
		wg.Add(1)
		go f(objCreate, objRemove, doneCh, &wg)
	}

	for notificationInfo := range self.client.ListenBucketNotification(bucketName, "", "", []string{
		"s3:ObjectCreated:*",
		"s3:ObjectRemoved:*",
	}, doneCh) {
		if notificationInfo.Err != nil {
			log.Println(notificationInfo.Err)
		}
		log.Println(notificationInfo)

		for _, val := range notificationInfo.Records {
			if strings.HasPrefix(val.EventName, "s3:ObjectCreated:") {
				objCreate <- StorageNotificationEvent(val)
			} else if strings.HasPrefix(val.EventName, "s3:ObjectRemoved:") {
				objRemove <- StorageNotificationEvent(val)
			}

		}
	}
	wg.Wait()
}
