package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/maxkrivich/ProjectX/api"
	"github.com/maxkrivich/ProjectX/api/storage"
	"github.com/maxkrivich/ProjectX/configs"
	"github.com/maxkrivich/ProjectX/models"
)

type FileStorageService struct {
	config  *configs.Config
	api     *api.APIControllers
	db      *models.DB
	storage *storage.StorageClient
}

func NewFileStorageService(conf *configs.Config) *FileStorageService {
	fss := FileStorageService{}
	fss.config = conf
	fss.storage = storage.DialStorage()
	fss.db = models.GetDB()
	fss.api = api.NewAPIContorller(fss.storage, fss.config)
	return &fss
}

func (self *FileStorageService) Run() {
	go self.storage.ListenBucketNotification(self.config.FileBucketName, models.HandleFileEvents)
	err := self.api.ListenAndServe(fmt.Sprintf("%s:%s", self.config.ServerConfig.Host, self.config.ServerConfig.Port))
	if err != nil {
		log.Fatal(err)
	}
}
func listenInterrupt() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case <-sig:
			log.Println("Terminating...")
			return
		}
	}
}

func main() {
	flag.Parse()
	api := NewFileStorageService(configs.GetConfig())
	defer api.db.Close()
	go api.Run()
	listenInterrupt()
}
