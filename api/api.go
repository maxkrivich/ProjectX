package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxkrivich/ProjectX/api/storage"
	"github.com/maxkrivich/ProjectX/configs"
)

// API messages
var (
	errors = map[int]ApiError{
		http.StatusBadRequest:          ApiError{StatusCode: http.StatusBadRequest, ErrorMessage: "Bad request"},
		http.StatusNotFound:            ApiError{StatusCode: http.StatusNotFound, ErrorMessage: "Not found"},
		http.StatusInternalServerError: ApiError{StatusCode: http.StatusInternalServerError, ErrorMessage: "=("},
	}
)

type ApiError struct {
	StatusCode   int
	ErrorMessage string
}

type APIControllers struct {
	storage *storage.StorageClient
	configs *configs.Config
	router  *gin.Engine
}

func NewAPIContorller(storage *storage.StorageClient, config *configs.Config) *APIControllers {
	api := APIControllers{
		storage: storage,
		configs: config,
		router:  gin.Default(),
	}
	api.initRouters()
	return &api
}

func (self *APIControllers) initRouters() {
	// init 'run' routers
	self.router.GET("/run", self.GetAllRun)
	self.router.GET("/run/:id", self.GetRun)
	self.router.POST("/run", self.CreateRun)
	self.router.PUT("/run/:id", self.UpdateRun)
	self.router.DELETE("/run/:id", self.DeleteRun)

	// init 'file' routers
	self.router.GET("/file/download", self.FileDownload)
	self.router.POST("/file/upload", self.FileUpload)
	self.router.DELETE("/file/delete", self.FileDelete)
}

func (self *APIControllers) ListenAndServe(addr string) error {
	return self.router.Run()
}
