package dao

import (
	"FunPDF/internal/dto"
	"FunPDF/internal/entity"
	"context"

	"gorm.io/gorm"
)

type FileDAO struct{}

func NewFileDAO() *FileDAO {
	return &FileDAO{}
}

/* Helper */

// CheckUUIO check one ID is exist
func (d *FileDAO) CheckUUIO(ctx context.Context, ID string, db *gorm.DB) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(&entity.File{}).
		Where("id = ?", ID).
		Count(&count).Error
	return count > 0, err
}

/* API */

// ListFiles List all files add to APP
func (d *FileDAO) ListFiles(ctx context.Context, db *gorm.DB) ([]entity.File, error) {
	fileList := make([]entity.File, 0)
	err := db.WithContext(ctx).Model(&entity.File{}).Find(&fileList).Error
	return fileList, err
}

// UpdateFileByID update file by file ID
func (d *FileDAO) UpdateFileByID(ctx context.Context, fileID string, db *gorm.DB, req *dto.UpdateFileRequest) (int64, error) {
	result := db.WithContext(ctx).Model(&entity.File{}).
		Where("id = ? AND revision = ?", fileID, req.ExpectedRevision).
		Updates(map[string]any{
			"revision":  req.Revision,
			"name":      req.Name,
			"mime_type": req.MimeType,
			"status":    req.Status,
		})
	return result.RowsAffected, result.Error
}

// GetFileByID get the file by file ID
func (d *FileDAO) GetFileByID(ctx context.Context, fileID string, db *gorm.DB) (*entity.File, error) {
	var file entity.File
	err := db.WithContext(ctx).First(&file, "id = ?", fileID).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// AdvanceRevision update `revision` only
func (d *FileDAO) AdvanceRevision(ctx context.Context, fileID string, db *gorm.DB, expectedRevision, nextRevision int64) (int64, error) {
	result := db.WithContext(ctx).Model(&entity.File{}).
		Where("id = ? AND revision = ?", fileID, expectedRevision).
		Updates(map[string]any{
			"revision": nextRevision,
		})
	return result.RowsAffected, result.Error
}

func (d *FileDAO) UploadFile(ctx context.Context, file *entity.File, db *gorm.DB) (int64, error) {
	result := db.WithContext(ctx).Model(&entity.File{}).
		Create(file)
	return result.RowsAffected, result.Error
}

// DeleteFile delete the file by ID
func (d *FileDAO) DeleteFile(ctx context.Context, fileID string, db *gorm.DB) (int64, error) {
	var affected int64

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&entity.File{}).
			Where("id = ?", fileID).
			Delete(&entity.File{})
		affected = result.RowsAffected
		return result.Error
	})

	return affected, err
}
