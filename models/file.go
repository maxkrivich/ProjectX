package models

import (
	"github.com/satori/go.uuid"
)

type File struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	UUID     string `json:"-"`
	FileName string `json:"file_name"`
	RunID    uint   `json:"run_id"`
	FileID   uint   `json:"file_id"`
	Uploaded bool   `json:"-"`
}

func NewFile(fileName string, runId, fileId uint) *File {
	var file File
	uName := uuid.NewV4()
	file.UUID = uName.String()
	file.RunID = runId
	file.FileID = fileId
	file.Uploaded = false
	file.FileName = fileName
	return &file
}

func (self *File) FindByID(id uint) bool {
	return dbConn.db.First(&self, id).RecordNotFound()
}

func (self *File) FindByRunIdAndFileId(fileId, runId uint) bool {
	return dbConn.db.Where("file_id = ? AND run_id = ?", fileId, runId).First(&self).RecordNotFound()
}

func (self *File) FindByUUID(ulid string) bool {
	return dbConn.db.Where("uuid = ?", ulid).First(&self).RecordNotFound()
}

func (self *File) SaveToDB() error {
	return dbConn.db.Save(&self).Error
}

func (self *File) DeleteFromDB() error {
	return dbConn.db.Delete(&self).Error
}
