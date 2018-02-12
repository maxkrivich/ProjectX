package api

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/maxkrivich/ProjectX/models"

	"github.com/gin-gonic/gin"
)

func (self *APIControllers) FileUpload(c *gin.Context) {
	runId := c.Query("runID")
	fileId := c.Query("fileID")
	fileName := c.Query("name")

	rid, err := strconv.Atoi(runId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	fid, err := strconv.ParseUint(fileId, 10, 32)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	if len(fileName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}

	var run models.Run
	if run.FindByID(uint(rid)) {
		c.JSON(http.StatusNotFound, gin.H{"error": errors[http.StatusNotFound].ErrorMessage})
		return
	}
	var file models.File
	if file.FindByRunIdAndFileId(uint(fid), uint(rid)) {
		file = *models.NewFile(fileName, uint(rid), uint(fid))
	}
	if file.FileName != fileName {
		file.FileName = fileName
	}
	file.Uploaded = false

	if err = file.SaveToDB(); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	presignedURL, err := self.storage.PresignedPutObject(self.configs.FileBucketName, file.UUID, self.configs.PresignedUrlExpires)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors[http.StatusInternalServerError].ErrorMessage})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"upload_url": presignedURL.String()})
}

func (self *APIControllers) FileDownload(c *gin.Context) {
	runId := c.Query("runID")
	fileId := c.Query("fileID")

	rid, err := strconv.Atoi(runId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}

	fid, err := strconv.Atoi(fileId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}

	var file models.File
	if file.FindByRunIdAndFileId(uint(fid), uint(rid)) {
		c.JSON(http.StatusNotFound, gin.H{"error": errors[http.StatusNotFound].ErrorMessage})
		return
	}
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileName))
	presignedURL, err := self.storage.PresignedGetObject(self.configs.FileBucketName, file.UUID, self.configs.PresignedUrlExpires, reqParams)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors[http.StatusInternalServerError].ErrorMessage})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"download_url": presignedURL.String(), "file_name": file.FileName})
}

func (self *APIControllers) FileDelete(c *gin.Context) {
	runId := c.Query("runID")
	fileId := c.Query("fileID")

	rid, err := strconv.Atoi(runId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}

	fid, err := strconv.Atoi(fileId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}

	var file models.File
	if file.FindByRunIdAndFileId(uint(fid), uint(rid)) {
		c.JSON(http.StatusNotFound, gin.H{"error": errors[http.StatusNotFound].ErrorMessage})
		return
	}

	err = self.storage.RemoveObject(self.configs.FileBucketName, file.UUID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors[http.StatusInternalServerError].ErrorMessage})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Ok"})
}
