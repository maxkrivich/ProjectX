package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/maxkrivich/ProjectX/models"

	"github.com/gin-gonic/gin"
)

func (self *APIControllers) CreateRun(c *gin.Context) {
	var run models.Run
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(bodyBytes, &run); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	if !run.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	run.CreatedAt = time.Now()
	if err := run.SaveToDB(); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	c.JSON(http.StatusCreated, run)
}

func (self *APIControllers) GetRun(c *gin.Context) {
	id, err := self.getId(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	var run models.Run

	if run.FindByID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
	} else {
		c.JSON(http.StatusOK, run)
	}
}

func (self *APIControllers) GetAllRun(c *gin.Context) {
	limitQuery := c.DefaultQuery("perPage", "25")
	pageQuery := c.DefaultQuery("page", "1")
	var err error
	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
	}
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
	}
	pag, err := models.PaginateRuns(models.NewPagination(limit, page))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors[http.StatusInternalServerError].ErrorMessage})
		return
	}
	c.JSON(http.StatusOK, pag)
}

func (self *APIControllers) UpdateRun(c *gin.Context) {
	id, err := self.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	var run models.Run
	var existing models.Run
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(bodyBytes, &run); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	if !run.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	run.ID = uint(id)

	if existing.FindByID(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": errors[http.StatusNotFound].ErrorMessage})
	} else {
		run.CreatedAt = existing.CreatedAt
		run.UpdatedAt = time.Now()
		if err := run.SaveToDB(); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors[http.StatusInternalServerError].ErrorMessage})
			return
		}
		c.JSON(http.StatusOK, run)
	}
}

func (self *APIControllers) DeleteRun(c *gin.Context) {
	id, err := self.getId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors[http.StatusBadRequest].ErrorMessage})
		return
	}
	var run models.Run

	if run.FindByID(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": errors[http.StatusNotFound].ErrorMessage})
	} else {
		if err := run.DeleteFromDB().Error; err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors[http.StatusInternalServerError].ErrorMessage})
			return
		}
		c.Data(http.StatusNoContent, "application/json", make([]byte, 0))
	}
}

func (self *APIControllers) getId(c *gin.Context) (uint, error) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return uint(id), nil
}
