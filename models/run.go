package models

import (
	"time"
)

type Run struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	deletedAt *time.Time `sql:"index"`
	Name      string     `gorm:"size:255; not null" json:"name" binding:"required"`
}
