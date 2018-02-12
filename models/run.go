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

type RunPagination struct {
	Items []Run `json:"items"`
	Pagination
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

func NewRun(name string) *Run {
	var run Run
	run.CreatedAt = time.Now()
	run.UpdatedAt = run.CreatedAt
	run.Name = name
	return &run
}

func (self *Run) FindByID(id uint) bool {
	return dbConn.db.First(&self, id).RecordNotFound()
}

func (self *Run) SaveToDB() error {
	return dbConn.db.Save(&self).Error
}

func (self *Run) DeleteFromDB() error {
	return dbConn.db.Delete(&self).Error
}

func PaginateRuns(pag *Pagination) (*RunPagination, error) {
	var runs []Run
	tmp := dbConn.db.Offset(pag.Limit * (pag.Page - 1)).Order("id asc").Find(&runs)
	if tmp.Error != nil {
		return nil, tmp.Error
	}
	pag.Count = len(runs)
	tmp.Limit(pag.Limit).Find(&runs)
	if tmp.Error != nil {
		return nil, tmp.Error
	}
	pag.HasNext = pag.Count > pag.Limit*pag.Page
	return &RunPagination{Items: runs, Pagination: *pag}, nil
}
