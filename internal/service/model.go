package service

import (
	"FunPDF/internal/dao"
	"FunPDF/internal/dto"
	"FunPDF/internal/entity"
	"FunPDF/internal/entity/models"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type ModelService struct {
	httpClient *http.Client
}

func NewModelService() *ModelService {
	return &ModelService{
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext:           (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
				ResponseHeaderTimeout: 30 * time.Second,
			},
		},
	}
}

func (s *ModelService) ChatToModelStreamWithSender(ctx context.Context, providerID, modelName, modelID string, messages []models.Message, modelCfg models.ModelConfig, chatCfg models.ChatConfig, sender func(*string, *string) error) error {
	resp, err := s.ChatToModel(ctx, providerID, modelName, modelID, messages, modelCfg, chatCfg, sender)
	if err != nil {
		return err
	}
	if resp == nil {
		return nil
	}
	done := "[DONE]"
	return sender(&done, nil)
}

func (s *ModelService) ChatToModel(ctx context.Context, providerID, modelName, modelID string, messages []models.Message, modelCfg models.ModelConfig, chatCfg models.ChatConfig, sender ...func(*string, *string) error) (*dto.ChatResponse, error) {
	providerID = strings.TrimSpace(providerID)
	modelName = strings.TrimSpace(modelName)
	if modelName == "" {
		modelName = strings.TrimSpace(modelID)
	}
	if providerID == "" {
		return nil, errors.New("provider id is empty")
	}
	if modelName == "" {
		return nil, errors.New("model name is empty")
	}

	var provider entity.Provider
	if err := dao.DB.WithContext(ctx).Where("id = ?", providerID).First(&provider).Error; err != nil {
		return nil, err
	}

	apiKey := strings.TrimSpace(provider.APIKey)
	modelCfg.APIKey = &apiKey
	var streamSender func(*string, *string) error
	if len(sender) > 0 {
		streamSender = sender[0]
	}

	switch strings.ToLower(strings.TrimSpace(provider.Name)) {
	case "deepseek":
		return (&models.DeepSeekModel{
			BaseModel: models.BaseModel{
				BaseURL:    strings.TrimRight(provider.BaseURL, "/"),
				URLSuffix:  strings.TrimLeft(provider.URLSuffix["chat"], "/"),
				HTTPClient: s.httpClient,
			},
		}).Chat(ctx, &modelCfg, &chatCfg, messages, modelName, streamSender)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider.Name)
	}
}

func (s *ModelService) ListSupportedModels(ctx context.Context, providerID string, modelCfg models.ModelConfig) (*[]dto.ListModelsResponse, error) {
	providerID = strings.TrimSpace(providerID)
	if providerID == "" {
		return nil, errors.New("provider id is empty")
	}

	var provider entity.Provider
	if err := dao.DB.WithContext(ctx).Where("id = ?", providerID).First(&provider).Error; err != nil {
		return nil, err
	}

	apiKey := strings.TrimSpace(provider.APIKey)
	modelCfg.APIKey = &apiKey

	switch strings.ToLower(strings.TrimSpace(provider.Name)) {
	case "deepseek":
		return (&models.DeepSeekModel{
			BaseModel: models.BaseModel{
				BaseURL:    strings.TrimRight(provider.BaseURL, "/"),
				URLSuffix:  strings.TrimLeft(provider.URLSuffix["models"], "/"),
				HTTPClient: s.httpClient,
			},
		}).ListModels(ctx, &modelCfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider.Name)
	}
}
