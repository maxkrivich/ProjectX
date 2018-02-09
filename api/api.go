package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/maxkrivich/ProjectX/configs"
	"github.com/maxkrivich/ProjectX/models"

	"github.com/minio/minio-go"
	"github.com/oklog/ulid"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type APIControllers struct {
	db   *gorm.DB
	mc   *minio.Client
	conf *configs.Config
}

type Pagination struct {
	Limit   int  `json:"perPage"`
	Page    int  `json:"page"`
	Count   int  `json:"count"`
	HasNext bool `json:"hasNext"`
}

type RunPagination struct {
	Items []models.Run `json:"items"`
	Pagination
}

func NewAPIController(db *gorm.DB, mc *minio.Client, conf *configs.Config) *APIControllers {
	return &APIControllers{db: db, mc: mc, conf: conf}
}

func (rc *APIControllers) CreateRun(c *gin.Context) {
	var run models.Run
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(bodyBytes, &run); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if !run.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	run.CreatedAt = time.Now()
	rc.db.Save(&run)
	c.JSON(http.StatusCreated, run)
}

func (rc *APIControllers) GetRun(c *gin.Context) {
	id, err := rc.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var run models.Run

	if rc.db.First(&run, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	} else {
		c.JSON(http.StatusOK, run)
	}
}

func (rc *APIControllers) GetAllRun(c *gin.Context) {
	var pag Pagination
	limitQuery := c.DefaultQuery("perPage", "25")
	pageQuery := c.DefaultQuery("page", "1")
	var runs []models.Run
	var err error
	pag.Limit, err = strconv.Atoi(limitQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	}
	pag.Limit = int(math.Max(1, math.Min(10000, float64(pag.Limit))))

	pag.Page, err = strconv.Atoi(pageQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	}
	pag.Page = int(math.Max(1, float64(pag.Page)))

	tmp := rc.db.Offset(pag.Limit * (pag.Page - 1)).Order("id asc").Find(&runs)
	pag.Count = len(runs)
	tmp.Limit(pag.Limit).Find(&runs)
	pag.HasNext = pag.Count > pag.Limit*pag.Page
	c.JSON(http.StatusOK, RunPagination{Items: runs, Pagination: pag})
}

func (rc *APIControllers) UpdateRun(c *gin.Context) {
	id, err := rc.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var run models.Run
	var existing models.Run
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(bodyBytes, &run); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if !run.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	run.ID = uint(id)

	if rc.db.First(&existing, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	} else {
		run.CreatedAt = existing.CreatedAt
		run.UpdatedAt = time.Now()
		rc.db.Save(&run)
		c.JSON(http.StatusOK, run)
	}
}

func (rc *APIControllers) DeleteRun(c *gin.Context) {
	id, err := rc.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var run models.Run

	if rc.db.First(&run, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	} else {
		rc.db.Delete(&run)
		c.Data(http.StatusNoContent, "application/json", make([]byte, 0))
	}
}

func (rc *APIControllers) getId(c *gin.Context) (uint, error) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return uint(id), nil
}

func (rc *APIControllers) FileUpload(c *gin.Context) {
	runId := c.Query("runID")
	fileId := c.Query("fileID")
	fileName := c.Query("name")

	rid, err := strconv.Atoi(runId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	fid, err := strconv.ParseUint(fileId, 10, 32)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if len(fileName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file name is empty"})
		return
	}

	var run models.Run
	if rc.db.First(&run, rid).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var file models.File

	if rc.db.Model(&file).Where("file_id = ? AND run_id = ?", fid, rid).First(&file).RecordNotFound() {
		t := time.Unix(1000000, 0)
		entropy := rand.New(rand.NewSource(t.UnixNano()))
		uName := ulid.MustNew(ulid.Timestamp(t), entropy)
		file.ULID = uName.String()
		file.RunID = run.ID // Dirty read if run was deleted
		file.FileID = uint(fid)
		file.Uploaded = false
		file.FileName = fileName
	}
	if file.FileName != fileName {
		file.FileName = fileName
	}
	file.Uploaded = false
	if err = rc.db.Save(&file).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	presignedURL, err := rc.mc.PresignedPutObject(rc.conf.FileBucketName, file.ULID, rc.conf.PresignedUrlExpires)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "=("})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"upload_url": presignedURL.String()})
}

func (rc *APIControllers) FileDownload(c *gin.Context) {
	runId := c.Query("runID")
	fileId := c.Query("fileID")

	rid, err := strconv.Atoi(runId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	fid, err := strconv.Atoi(fileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	var file models.File
	if rc.db.Model(&file).Where("file_id = ? AND run_id = ?", fid, rid).First(&file).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileName))
	presignedURL, err := rc.mc.PresignedGetObject(rc.conf.FileBucketName, file.ULID, rc.conf.PresignedUrlExpires, reqParams)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "=("})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"download_url": presignedURL.String(), "file_name": file.FileName})
}

func (rc *APIControllers) FileDelete(c *gin.Context) {

}
