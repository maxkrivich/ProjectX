package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/maxkrivich/ProjectX/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RunControllers struct {
	db *gorm.DB
}

func NewRunControllers(db *gorm.DB) *RunControllers {
	return &RunControllers{db: db}
}

func (rc *RunControllers) CreateRun(c *gin.Context) {
	var run models.Run
	if c.Bind(&run) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "problem decoding id sent"})
		return
	}
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
	var runs []models.Run
	rc.db.Order("created_at desc").Find(&runs)
	c.JSON(http.StatusOK, runs)
}

func (rc *RunControllers) UpdateRun(c *gin.Context) {
	id, err := rc.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var run models.Run
	var existing models.Run
	if c.Bind(&run) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "problem decoding id sent"})
	}
	run.ID = uint(id)

	if rc.db.First(&existing, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	} else {
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
