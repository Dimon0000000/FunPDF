package handler

import (
	"FunPDF/internal/dto"
	"FunPDF/internal/service"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileSvr *service.FileService
}

func NewFileHandler() *FileHandler {
	return &FileHandler{fileSvr: service.NewFileService()}
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	fileList, err := h.fileSvr.ListFiles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": fileList, "msg": "success"})
}

// SaveFile saves the edited JSON state to local storage.
func (h *FileHandler) SaveFile(c *gin.Context) {
	fileID := strings.TrimSpace(c.Param("file_id"))
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "file id is empty"})
		return
	}

	var req dto.SaveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	ok, err := h.fileSvr.SaveFile(c.Request.Context(), fileID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "revision") {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"code": status, "msg": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "save failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success"})
}

// UploadFile stores the PDF and initial editor state on first Ctrl+S.
func (h *FileHandler) UploadFile(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 200<<20)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	stateText := strings.TrimSpace(c.PostForm("editor_state"))
	if stateText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "editor_state is empty"})
		return
	}
	if !json.Valid([]byte(stateText)) {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "editor_state is invalid JSON"})
		return
	}

	source, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	defer source.Close()

	req := &dto.UploadFileRequest{
		FileName:    fileHeader.Filename,
		MimeType:    fileHeader.Header.Get("Content-Type"),
		FileSize:    fileHeader.Size,
		EditorState: json.RawMessage(stateText),
	}
	file, err := h.fileSvr.UploadFile(c.Request.Context(), req, source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "msg": "success", "data": file})
}

// AlertFile update the file metadata
func (h *FileHandler) AlertFile(c *gin.Context) {

}

// DeleteFile delete the file
func (h *FileHandler) DeleteFile(c *gin.Context) {
	fileID := strings.TrimSpace(c.Param("file_id"))
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "file id is empty"})
		return
	}
	affected, err := h.fileSvr.DeleteFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "msg": "file id is not exist"})
		return
	}
	c.Status(http.StatusNoContent)
}
