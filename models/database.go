package models

import (
	"log"

	"github.com/maxkrivich/ProjectX/configs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DB struct {
	db *gorm.DB
}

var dbConn *DB

func init() {
	dbConn = GetDB()
}

func GetDB() *DB {
	if dbConn != nil {
		return dbConn
	}
	dbConn, err := newDB()
	if err != nil {
		log.Fatal(err)
	}
	return dbConn
}

func newDB() (*DB, error) {
	config := configs.GetConfig().DBConfig
	var err error
	var dbC DB
	db, err := gorm.Open(config.Dialect, config.GetConnectionString())
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	dbC.db = db
	dbC.initSchema()
	return &dbC, nil
}

func (self *DB) initSchema() {
	if !self.db.HasTable(&Run{}) {
		self.db.CreateTable(&Run{})
	}
	if !self.db.HasTable(&File{}) {
		self.db.CreateTable(&File{})
	}
	self.db.Model(&File{}).AddForeignKey("run_id", "runs(id)", "RESTRICT", "RESTRICT")
	self.db.AutoMigrate(&File{}, &Run{})
}

func (self *DB) Close() {
	self.db.Close()
}
