package models

import (
	"time"
)

type File struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	UUID      string     `json:"-"`
	FileName  string     `json:"-"`
	RunID     uint
	FileID    uint `gorm:"not null;unique"`
}
