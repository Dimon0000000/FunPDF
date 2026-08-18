package internal

import (
	"FunPDF/internal/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	fileHandler  *handler.FileHandler
	albumHandler *handler.AlbumHandler
}

// NewRouter create a new router
func NewRouter(fileHandler *handler.FileHandler, albumHandler *handler.AlbumHandler) *Router {
	return &Router{
		fileHandler:  fileHandler,
		albumHandler: albumHandler,
	}
}

// Setup register all API routes
func (r *Router) Setup(e *gin.Engine) {
	api := e.Group("/api")
	{
		file := api.Group("/files")
		{
			file.GET("", r.fileHandler.ListFiles)
			file.POST("", r.fileHandler.UploadFile)
			file.PUT("/:file_id", r.fileHandler.AlertFile) // TODO finish this API
			file.DELETE("/:file_id", r.fileHandler.DeleteFile)
			file.PATCH("/:file_id/state", r.fileHandler.SaveFile)
		}

		album := api.Group("/album")
		{
			album.GET("", r.albumHandler.ListAlbums)
			album.POST("", r.albumHandler.CreateAlbum)

			album.GET("/:album_id", r.albumHandler.ListAlbumFiles)
			album.PUT("/:album_id", r.albumHandler.UpdateAlbum)
			album.DELETE("/:album_id", r.albumHandler.DeleteAlbum)

			album.POST(":album_id/files", r.albumHandler.UploadFilesToAlbum)
			album.DELETE("/:album_id/files", r.albumHandler.DeleteFilesFromAlbum)
			// album.DELETE("/:album_id/files/delete", r.albumHandler.HardDeleteFilesFromAlbum) TODO: will be added in v0.2.x
		}
	}
}
