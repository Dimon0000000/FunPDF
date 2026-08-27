package service

import (
	"FunPDF/internal/dao"
	"FunPDF/internal/dto"
	"FunPDF/internal/entity"
	"FunPDF/internal/entity/models"
	"context"
	"errors"
	"fmt"
	"strings"
)

type ModelService struct {
}

func NewModelService() *ModelService {
	return &ModelService{}
}

func (s *ModelService) ChatToModelStreamWithSender(ctx context.Context, providerID, modelName, modelID string, messages []models.Message, modelCfg models.ModelConfig, chatCfg models.ChatConfig, sender func(*string, *string) error) error {
	resp, err := s.ChatToModel(ctx, providerID, modelName, modelID, messages, modelCfg, chatCfg)
	if err != nil {
		return err
	}
	if resp.ReasonContent != nil && *resp.ReasonContent != "" {
		if err := sender(nil, resp.ReasonContent); err != nil {
			return err
		}
	}
	if resp.Answer != nil && *resp.Answer != "" {
		if err := sender(resp.Answer, nil); err != nil {
			return err
		}
	}
	done := "[DONE]"
	return sender(&done, nil)
}

func (s *ModelService) ChatToModel(ctx context.Context, providerID, modelName, modelID string, messages []models.Message, modelCfg models.ModelConfig, chatCfg models.ChatConfig) (*dto.ChatResponse, error) {
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

	switch strings.ToLower(strings.TrimSpace(provider.Name)) {
	case "deepseek":
		return (&models.DeepSeekModel{
			BaseModel: models.BaseModel{
				BaseURL:   strings.TrimRight(provider.BaseURL, "/"),
				URLSuffix: strings.TrimLeft(provider.URLSuffix["chat"], "/"),
			},
		}).Chat(ctx, &modelCfg, &chatCfg, messages, modelName)
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
				BaseURL:   strings.TrimRight(provider.BaseURL, "/"),
				URLSuffix: strings.TrimLeft(provider.URLSuffix["models"], "/"),
			},
		}).ListModels(ctx, &modelCfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider.Name)
	}
}
