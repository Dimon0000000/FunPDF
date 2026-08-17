package service

import (
	"FunPDF/internal/common"
	"FunPDF/internal/dao"
	"FunPDF/internal/dto"
	"FunPDF/internal/entity"
	"context"
	"fmt"
	"strings"
)

type AlbumService struct {
	albumDAO *dao.AlbumDAO
}

func NewAlbumService() *AlbumService {
	return &AlbumService{albumDAO: dao.NewAlbumDAO()}
}

// ListAlbums list all albums
func (s *AlbumService) ListAlbums(ctx context.Context) ([]*entity.Album, error) {
	albums, err := s.albumDAO.ListAlbums(ctx, dao.DB)
	if err != nil {
		return nil, err
	}
	return albums, nil
}

// CreateAlbum create an album
func (s *AlbumService) CreateAlbum(ctx context.Context, req *dto.CreateAlbumReq) (*entity.Album, error) {
	albumID := common.GenerateUUIDv7()

	// check thumbnail size
	if err := ValidateBase64ImageSize(req.Thumbnail); err != nil {
		return nil, err
	}

	album := &entity.Album{
		ID:          albumID,
		Name:        strings.TrimSpace(req.Name),
		Thumbnail:   req.Thumbnail,
		Description: req.Description,
	}

	affected, err := s.albumDAO.CreateAlbum(ctx, dao.DB, album)
	if err != nil || affected == 0 {
		return nil, fmt.Errorf("create album failed: %w", err)
	}
	return album, nil
}

// ListAlbumFiles list all files under the album
func (s *AlbumService) ListAlbumFiles(ctx context.Context, albumID string) ([]*entity.File, error) {
	result, err := s.albumDAO.ListAlbumFiles(ctx, dao.DB, albumID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
