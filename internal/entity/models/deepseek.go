package models

import (
	"FunPDF/internal/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DeepSeekModel struct {
	BaseModel
}

func (d *DeepSeekModel) Chat(ctx context.Context, modelCfg *ModelConfig, chatCfg *ChatConfig, messages []Message, modelName string) (*dto.ChatResponse, error) {
	url := fmt.Sprintf("%s/%s", d.BaseURL, d.URLSuffix)

	isStream := chatCfg != nil && chatCfg.Stream != nil && *chatCfg.Stream
	reqBody := map[string]any{
		"model":    modelName,
		"messages": messages,
	}
	if chatCfg != nil {
		implementDeepSeekChatConfig(chatCfg, reqBody)
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	toCtx := ctx
	var cancel context.CancelFunc
	if d.httpClient != nil && d.httpClient.Timeout > 0 {
		toCtx, cancel = context.WithTimeout(ctx, d.httpClient.Timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(toCtx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *modelCfg.APIKey))

	if isStream {
		return doStreamChat(req)
	}
	return doNoneStreamChat(req)

}

func (d *DeepSeekModel) ListModels(ctx context.Context, modelCfg *ModelConfig) (*[]dto.ListModelsResponse, error) {
	url := fmt.Sprintf("%s/%s", d.BaseURL, d.URLSuffix)

	toCtx := ctx
	var cancel context.CancelFunc
	if d.httpClient != nil && d.httpClient.Timeout > 0 {
		toCtx, cancel = context.WithTimeout(ctx, d.httpClient.Timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(toCtx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *modelCfg.APIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]any
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	modelList, ok := data["data"].([]any)
	if !ok {
		modelList, ok = data["models"].([]any)
	}
	if !ok {
		return nil, fmt.Errorf(`"models" is not a list`)
	}

	var models []dto.ListModelsResponse
	for _, v := range modelList {
		item, ok := v.(map[string]any)
		if !ok {
			continue
		}
		modelName, ok := item["id"].(string)
		if !ok {
			continue
		}
		models = append(models, dto.ListModelsResponse{Name: modelName})
	}

	return &models, nil
}

func (d *DeepSeekModel) Name() string {
	return "DeepSeek"
}
