package dao

import (
	"FunPDF/internal/entity"
	"context"

	"gorm.io/gorm"
)

type AlbumDAO struct {
}

func NewAlbumDAO() *AlbumDAO {
	return &AlbumDAO{}
}

// ListAlbums list all albums
func (d *AlbumDAO) ListAlbums(ctx context.Context, db *gorm.DB) ([]*entity.Album, error) {
	var albums []*entity.Album
	err := db.WithContext(ctx).Model(&entity.Album{}).Find(&albums).Error
	if err != nil {
		return nil, err
	}
	return albums, nil
}

// CreateAlbum create an album
func (d *AlbumDAO) CreateAlbum(ctx context.Context, db *gorm.DB, album *entity.Album) (int64, error) {
	result := db.WithContext(ctx).Model(&entity.Album{}).
		Create(&album)
	return result.RowsAffected, result.Error
}

// ListAlbumFiles list all files under the album
func (d *AlbumDAO) ListAlbumFiles(ctx context.Context, db *gorm.DB, albumID string) ([]*entity.File, error) {
	var files []*entity.File
	err := db.WithContext(ctx).Model(&entity.AlbumFile{}).
		Where("album_id = ?", albumID).Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}
