package models

import (
	"FunPDF/internal/dto"
	"context"
	"net/http"
)

/* Interface */

type BaseModelInterface interface {
	Chat(ctx context.Context, modelCfg *ModelConfig, chatCfg *ChatConfig, messages []Message, modelName string, sender func(*string, *string) error) (*dto.ChatResponse, error)
	ListModels(ctx context.Context, modelCfg *ModelConfig) (*[]dto.ListModelsResponse, error)
	Name() string
}

/* Struct */

type BaseModel struct {
	BaseURL    string
	URLSuffix  string
	HTTPClient *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type ModelConfig struct {
	APIKey *string `json:"api_key"`
	Region *string `json:"region"`
}

type ChatConfig struct {
	Stream      *bool
	Thinking    *bool
	MaxTokens   *int
	Temperature *float64
	TopP        *float64
	DoSample    *bool
	Stop        *[]string
	Effort      *string
	Verbosity   *string
}
