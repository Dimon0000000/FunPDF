package dao

import (
	"FunPDF/internal/common"
	"FunPDF/internal/dto"
	"FunPDF/internal/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type ProviderDAO struct {
}

func NewProviderDAO() *ProviderDAO {
	return &ProviderDAO{}
}

// ListProviders list all providers that is already in DB
func (d *ProviderDAO) ListProviders(ctx context.Context, db *gorm.DB) ([]dto.ListProvidersResult, error) {
	var providers []dto.ListProvidersResult
	err := db.WithContext(ctx).Model(&entity.Provider{}).Find(&providers).Error
	if err != nil {
		return nil, err
	}
	return providers, nil
}

// CreateProvider create a provider
func (d *ProviderDAO) CreateProvider(ctx context.Context, db *gorm.DB, req *dto.CreateProviderRequest) error {
	id := common.GenerateUUIDv7()
	provider := entity.Provider{
		ID:        id,
		Name:      req.Name,
		BaseURL:   req.BaseURL,
		URLSuffix: req.URLSuffix,
		APIKey:    req.APIKey,
	}
	err := db.WithContext(ctx).Model(&entity.Provider{}).Create(&provider).Error
	return err
}

// UpdateProvider update provider
func (d *ProviderDAO) UpdateProvider(ctx context.Context, db *gorm.DB, req *dto.UpdateProviderRequest, providerID string) (int64, error) {
	result := db.WithContext(ctx).Model(&entity.Provider{}).
		Where("id = ?", providerID).
		Updates(*req)
	return result.RowsAffected, result.Error
}

// DeleteProvider delete provider and it's models
func (d *ProviderDAO) DeleteProvider(ctx context.Context, db *gorm.DB, providerID string) error {
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// delete provider
		err := tx.WithContext(ctx).Model(&entity.Provider{}).
			Where("id = ?", providerID).
			Delete(&entity.Provider{}).
			Error
		if err != nil {
			return fmt.Errorf("delete provider error: %s", err.Error())
		}

		// delete models that related to provider
		var models []entity.ProviderModel
		err = db.WithContext(ctx).Where("provider_id = ?", providerID).Find(&models).Error
		if err != nil {
			return fmt.Errorf("delete provider error: %s", err.Error())
		}

		for _, model := range models {
			err = db.WithContext(ctx).Model(&entity.Model{}).
				Where("id = ?", model.ModelID).Error
			if err != nil {
				return fmt.Errorf("delete provider error: %s", err.Error())
			}
		}

		// delete relationship
		err = db.WithContext(ctx).Model(&entity.ProviderModel{}).
			Where("provider_id = ?", providerID).
			Delete(&entity.ProviderModel{}).Error
		if err != nil {
			return fmt.Errorf("delete provider-model relationship error: %s", err.Error())
		}

		return nil
	})
	return err
}
