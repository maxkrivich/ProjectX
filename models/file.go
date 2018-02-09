package models

type File struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	ULID     string `json:"-"`
	FileName string `json:"file_name"`
	RunID    uint   `json:"run_id"`
	FileID   uint   `json:"file_id"`
	Uploaded bool   `json:"-"`
}
