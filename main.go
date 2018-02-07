package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/maxkrivich/ProjectX/api"
	"github.com/maxkrivich/ProjectX/configs"
	"github.com/maxkrivich/ProjectX/models"
	"github.com/minio/minio-go"

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
	rs.initRouters()

	if err := rs.initMinioClient(); err != nil {
		log.Fatal(err)
	}

	return &rs
}

func (rs *RunService) initMinioClient() (err error) {
	rs.mc, err = minio.New(rs.configs.Endpoint, rs.configs.AccessKeyID, rs.configs.SecretAccessKey, rs.configs.UseSSL)
	if err != nil {
		return err
	}
	bucketName := "files"

	if err := rs.mc.MakeBucket(bucketName, "/data"); err != nil {

		exists, err := rs.mc.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("Bucket '%s' is already exists\n", bucketName)
		} else {
			return err
		}
	}
	log.Printf("Successfully created '%s'\n", bucketName)
	return nil
}

func (rs *RunService) Run() error {
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
	db.AutoMigrate(&models.File{}, &models.Run{})
	db.Model(&models.Run{}).Related(&models.File{})
	return db, nil
}

func (rs *RunService) initRouters() {
	rc := api.NewAPIController(rs.db)
	rs.router.GET("/run", rc.GetAllRun)
	rs.router.GET("/run/:id", rc.GetRun)
	rs.router.POST("/run", rc.CreateRun)
	rs.router.PUT("/run/:id", rc.UpdateRun)
	rs.router.DELETE("/run/:id", rc.DeleteRun)

	rs.router.GET("file/download", rc.FileDownload)
	rs.router.POST("file/upload", rc.FileUpload)
	rs.router.DELETE("file/delete", rc.FileDelete)
}

// TODO
//func (rs *RunService) Migrate() error {
//	return nil
//}

func main() {
	flag.Parse()
	api := NewRunService(configs.GetConfig())
	defer api.db.Close()
	err := api.Run()
	if err != nil {
		log.Fatal(err)
	}
}
