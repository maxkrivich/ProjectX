package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/maxkrivich/ProjectX/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RunControllers struct {
	db *gorm.DB
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

func NewRunControllers(db *gorm.DB) *RunControllers {
	return &RunControllers{db: db}
}

func (rc *RunControllers) CreateRun(c *gin.Context) {
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

func (rc *RunControllers) GetRun(c *gin.Context) {
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

func (rc *RunControllers) GetAllRun(c *gin.Context) {
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

func (rc *RunControllers) UpdateRun(c *gin.Context) {
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

func (rc *RunControllers) DeleteRun(c *gin.Context) {
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

func (rc *RunControllers) getId(c *gin.Context) (uint, error) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return uint(id), nil
}
