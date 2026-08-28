package service

import (
	"FunPDF/internal/dao"
	"FunPDF/internal/dto"
	"FunPDF/internal/entity"
	"context"
	"errors"
	"strings"
)

type ProviderService struct {
	providerDAO *dao.ProviderDAO
}

func NewProviderService() *ProviderService {
	return &ProviderService{
		providerDAO: dao.NewProviderDAO(),
	}
}

// ListProviders list all providers
func (s *ProviderService) ListProviders(ctx context.Context) (*[]dto.ListProvidersResult, error) {
	providers, err := s.providerDAO.ListProviders(ctx, dao.DB)
	if err != nil {
		return nil, err
	}
	return &providers, nil
}

// CreateProvider create a provider
func (s *ProviderService) CreateProvider(ctx context.Context, req *dto.CreateProviderRequest) (*entity.Provider, error) {
	modelName := strings.TrimSpace(req.Name)
	if modelName == "" {
		return nil, ErrProviderNameRequired
	}
	req.Name = modelName

	req.BaseURL = strings.TrimSpace(req.BaseURL)
	req.APIKey = strings.TrimSpace(req.APIKey)
	chatSuffix := strings.TrimSpace(req.URLSuffix["chat"])
	modelsSuffix := strings.TrimSpace(req.URLSuffix["models"])
	if chatSuffix == "" || modelsSuffix == "" {
		return nil, ErrProviderURLSuffix
	}
	req.URLSuffix["chat"] = chatSuffix
	req.URLSuffix["models"] = modelsSuffix

	provider, err := s.providerDAO.CreateProvider(ctx, dao.DB, req)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// UpdateProvider update provider
func (s *ProviderService) UpdateProvider(ctx context.Context, req *dto.UpdateProviderRequest, providerID string) error {
	providerID = strings.TrimSpace(providerID)
	if providerID == "" {
		return ErrProviderIDRequired
	}

	req.BaseURL = strings.TrimSpace(req.BaseURL)
	req.APIKey = strings.TrimSpace(req.APIKey)
	chatSuffix := strings.TrimSpace(req.URLSuffix["chat"])
	modelsSuffix := strings.TrimSpace(req.URLSuffix["models"])
	if chatSuffix == "" || modelsSuffix == "" {
		return ErrProviderURLSuffix
	}
	req.URLSuffix["chat"] = chatSuffix
	req.URLSuffix["models"] = modelsSuffix

	affected, err := s.providerDAO.UpdateProvider(ctx, dao.DB, req, providerID)
	if affected != 1 {
		if affected < 1 {
			return ErrProviderNotFound
		}

		return errors.New("affected too many rows. please delete the provider and create a new one")
	}
	if err != nil {
		return err
	}

	return nil
}

// DeleteProvider delete provider and it's models
func (s *ProviderService) DeleteProvider(ctx context.Context, providerID string) error {
	providerID = strings.TrimSpace(providerID)
	if providerID == "" {
		return errors.New("provider id is empty")
	}

	if err := s.providerDAO.DeleteProvider(ctx, dao.DB, providerID); err != nil {
		return err
	}
	return nil
}
