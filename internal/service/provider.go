package service

import (
	"FunPDF/internal/dao"
	"FunPDF/internal/dto"
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
func (s *ProviderService) CreateProvider(ctx context.Context, req *dto.CreateProviderRequest) error {
	modelName := strings.TrimSpace(req.Name)
	if modelName == "" {
		return errors.New("model name is empty")
	}
	req.Name = modelName

	req.BaseURL = strings.TrimSpace(req.BaseURL)
	req.APIKey = strings.TrimSpace(req.APIKey)

	if err := s.providerDAO.CreateProvider(ctx, dao.DB, req); err != nil {
		return err
	}
	return nil
}

// UpdateProvider update provider
func (s *ProviderService) UpdateProvider(ctx context.Context, req *dto.UpdateProviderRequest, providerID string) error {
	providerID = strings.TrimSpace(providerID)
	if providerID == "" {
		return errors.New("provider id is empty")
	}

	req.BaseURL = strings.TrimSpace(req.BaseURL)
	req.APIKey = strings.TrimSpace(req.APIKey)

	affected, err := s.providerDAO.UpdateProvider(ctx, dao.DB, req, providerID)
	if affected != 1 {
		if affected < 1 {
			return errors.New("provider not found")
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
