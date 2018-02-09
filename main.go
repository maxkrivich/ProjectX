package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/maxkrivich/ProjectX/api"
	"github.com/maxkrivich/ProjectX/configs"
	"github.com/maxkrivich/ProjectX/models"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/policy"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type RunService struct {
	configs *configs.Config
	router  *gin.Engine
	db      *gorm.DB
	mc      *minio.Client
}

func NewRunService(conf *configs.Config) *RunService {
	rs := RunService{}
	rs.configs = conf
	db, err := rs.initDB()
	if err != nil {
		log.Fatal(err)
	}
	rs.db = db
	rs.router = gin.Default()

	if err := rs.initMinio(); err != nil {
		log.Fatal(err)
	}

	rs.initRouters()

	return &rs
}

func (rs *RunService) initMinio() (err error) {
	rs.mc, err = minio.New(rs.configs.Endpoint, rs.configs.AccessKeyID, rs.configs.SecretAccessKey, rs.configs.UseSSL)
	if err != nil {
		return err
	}
	// setup bucket
	if err := rs.mc.MakeBucket(rs.configs.FileBucketName, ""); err != nil {

		exists, err := rs.mc.BucketExists(rs.configs.FileBucketName)
		if err == nil && exists {
			log.Printf("Bucket '%s' is already exists\n", rs.configs.FileBucketName)
		} else {
			return err
		}
	}
	log.Printf("Successfully created '%s'\n", rs.configs.FileBucketName)

	// setup bucket policy
	err = rs.mc.SetBucketPolicy(rs.configs.FileBucketName, "", policy.BucketPolicyReadWrite)
	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}

func (rs *RunService) Run() error {
	go rs.listenBucketNotification()
	err := rs.router.Run(fmt.Sprintf("%s:%s", rs.configs.ServerConfig.Host, rs.configs.ServerConfig.Port))
	return err
}

func (rs *RunService) initDB() (*gorm.DB, error) {
	db, err := gorm.Open(rs.configs.Dialect, rs.configs.GetConnectionString())
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	if !db.HasTable(&models.Run{}) {
		db.CreateTable(&models.Run{})
	}
	if !db.HasTable(&models.File{}) {
		db.CreateTable(&models.File{})
	}
	// init schema
	db.Model(&models.File{}).AddForeignKey("run_id", "runs(id)", "RESTRICT", "RESTRICT")
	db.AutoMigrate(&models.File{}, &models.Run{})
	return db, nil
}

func (rs *RunService) initRouters() {
	rc := api.NewAPIController(rs.db, rs.mc, rs.configs)
	rs.router.GET("/run", rc.GetAllRun)
	rs.router.GET("/run/:id", rc.GetRun)
	rs.router.POST("/run", rc.CreateRun)
	rs.router.PUT("/run/:id", rc.UpdateRun)
	rs.router.DELETE("/run/:id", rc.DeleteRun)

	rs.router.GET("file/download", rc.FileDownload)
	rs.router.POST("file/upload", rc.FileUpload)
	rs.router.DELETE("file/delete", rc.FileDelete)
}

func (rs *RunService) listenBucketNotification() {
	doneCh := make(chan struct{})
	objCreate, objRemove := make(chan minio.NotificationEvent), make(chan minio.NotificationEvent)
	defer close(doneCh)
	defer close(objRemove)
	defer close(objCreate)

	go rs.handleEvents(objCreate, objRemove, doneCh)

	for notificationInfo := range rs.mc.ListenBucketNotification(rs.configs.FileBucketName, "", "", []string{
		"s3:ObjectCreated:*",
		"s3:ObjectRemoved:*",
	}, doneCh) {
		if notificationInfo.Err != nil {
			log.Println(notificationInfo.Err)
		}
		log.Println(notificationInfo)

		for _, val := range notificationInfo.Records {
			if strings.HasPrefix(val.EventName, "s3:ObjectCreated:") {
				objCreate <- val
			} else if strings.HasPrefix(val.EventName, "s3:ObjectRemoved:") {
				objRemove <- val
			}

		}
	}
}

func (rs *RunService) handleEvents(objCreate, objRemove chan minio.NotificationEvent, done chan struct{}) {
	for {
		select {
		case e := <-objCreate:
			var file models.File
			if !rs.db.Model(&file).Where("uuid = ?", e.S3.Object.Key).First(&file).RecordNotFound() {
				file.Uploaded = true
				rs.db.Save(&file)
			}
		case e := <-objRemove:
			var file models.File
			if !rs.db.Model(&file).Where("uuid = ?", e.S3.Object.Key).First(&file).RecordNotFound() {
				rs.db.Delete(&file)
			}
		case <-done:
			return
		}
	}
}

func main() {
	flag.Parse()
	api := NewRunService(configs.GetConfig())
	defer api.db.Close()
	err := api.Run()
	if err != nil {
		log.Fatal(err)
	}
}
