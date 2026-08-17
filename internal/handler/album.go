package handler

import (
	"FunPDF/internal/dto"
	"FunPDF/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	albumSvr *service.AlbumService
}

func NewAlbumHandler() *AlbumHandler {
	return &AlbumHandler{albumSvr: service.NewAlbumService()}
}

// ListAlbums list all albums
func (h *AlbumHandler) ListAlbums(c *gin.Context) {
	// TODO The list logic here will be changed in v0.3, keep this simple for now
	result, err := h.albumSvr.ListAlbums(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": result,
		"msg":  "success",
	})
}

// CreateAlbum create an album
func (h *AlbumHandler) CreateAlbum(c *gin.Context) {
	var req *dto.CreateAlbumReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}

	// create the album
	album, err := h.albumSvr.CreateAlbum(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": album,
		"msg":  "success",
	})
}

// ListAlbumFiles list all files under the album
func (h *AlbumHandler) ListAlbumFiles(c *gin.Context) {
	albumID := strings.TrimSpace(c.Param("album_id"))
	if albumID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "invalid album id",
		})
		return
	}

	result, err := h.albumSvr.ListAlbumFiles(c.Request.Context(), albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": result,
		"msg":  "success",
	})
}

// UpdateAlbum update the album
func (h *AlbumHandler) UpdateAlbum(c *gin.Context) {

}

// DeleteAlbum delete the album
func (h *AlbumHandler) DeleteAlbum(c *gin.Context) {

}

// UploadFiles upload a batch of files to album
func (h *AlbumHandler) UploadFiles(c *gin.Context) {

}

// DeleteFiles delete a batch of files from album
func (h *AlbumHandler) DeleteFiles(c *gin.Context) {

}
