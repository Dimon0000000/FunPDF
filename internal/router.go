package internal

import (
	"FunPDF/internal/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	fileHandler *handler.FileHandler
}

// NewRouter create a new router
func NewRouter(fileHandler *handler.FileHandler) *Router {
	return &Router{
		fileHandler: fileHandler,
	}
}

// Setup register all API routes
func (r *Router) Setup(e *gin.Engine) {
	api := e.Group("/api")
	{
		files := api.Group("/files")
		{
			files.GET("", r.fileHandler.ListFiles)
			files.POST("", r.fileHandler.UploadFile)
			files.PATCH("/:file_id/state", r.fileHandler.SaveFile)
			files.DELETE("/:file_id", r.fileHandler.DeleteFile)
		}
	}
}
