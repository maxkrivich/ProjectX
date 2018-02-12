package models

import (
	"sync"

	"github.com/maxkrivich/ProjectX/api/storage"
)

func HandleFileEvents(objCreate, objRemove <-chan storage.StorageNotificationEvent, done chan struct{}, group *sync.WaitGroup) {
	for {
		select {
		case e := <-objCreate:
			var file File
			if !file.FindByUUID(e.S3.Object.Key) {
				file.Uploaded = true
				file.SaveToDB()
			}
		case e := <-objRemove:
			var file File
			if !file.FindByUUID(e.S3.Object.Key) {
				file.DeleteFromDB()
			}
		case <-done:
			group.Done()
			return
		}
	}
}
