package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/maxkrivich/ProjectX/api"
	"github.com/maxkrivich/ProjectX/configs"
	"github.com/maxkrivich/ProjectX/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type RunService struct {
	configs *configs.Config
	router  *gin.Engine
	db      *gorm.DB
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
	return &rs
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
