package models

import (
	"time"
)

type Run struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	Name      string     `gorm:"size:255; not null" json:"name"`
}

func (r Run) Validate() bool {
	if len(r.Name) == 0 {
		return false
	}
	if r.ID != 0 || r.UpdatedAt != (time.Time{}) || r.CreatedAt != (time.Time{}) || r.DeletedAt != nil {
		return false
	}
	return true
}
