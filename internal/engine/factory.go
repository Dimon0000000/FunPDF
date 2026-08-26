package engine

import (
	"FunPDF/internal/dao"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type TranslatorFactory struct {
	translatorDAO *dao.TranslatorDAO
}

func NewTranslatorFactory() *TranslatorFactory {
	return &TranslatorFactory{
		translatorDAO: &dao.TranslatorDAO{},
	}
}

func NormalizeTranslatorName(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	switch normalized {
	case "baidu", "baidu-translator", "baidutranslator":
		return "baidu"
	case "deepl", "deepl-translator", "deep-l-translator", "deepltranslator":
		return "deepl"
	default:
		return normalized
	}
}

func (t *TranslatorFactory) GetTranslator(ctx context.Context, db *gorm.DB, translatorName, region string) (Translator, error) {
	translatorName = NormalizeTranslatorName(translatorName)
	if translatorName == "" {
		return nil, errors.New("translator name is required")
	}

	// get translator params from DB
	params, err := t.translatorDAO.GetTranslatorParams(ctx, db, translatorName)
	if err != nil {
		return nil, err
	}
	var param map[string]any
	err = json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}

	// create translator instance by name
	switch translatorName {
	case "baidu":
		apiKey, ok := param["api_key"].(string)
		if !ok {
			return nil, errors.New("api_key is invalid")
		}
		appID, ok := param["app_id"].(string)
		if !ok {
			return nil, errors.New("app_id is invalid")
		}
		url, err := translatorConfigURL("baidu", "default")
		if err != nil {
			return nil, err
		}
		return NewBaiduTranslator(apiKey, appID, url), nil
	case "deepl":
		apiKey, ok := param["api_key"].(string)
		if !ok {
			return nil, errors.New("api_key is invalid")
		}

		if region != "free" && region != "pro" {
			return nil, fmt.Errorf("region %s not supported for %s", region, translatorName)
		}
		url, err := translatorConfigURL("deepl", region)
		if err != nil {
			return nil, err
		}

		return NewDeeplTranslator(apiKey, url), nil
	default:
		return nil, fmt.Errorf("unsupported translator: %s", translatorName)
	}
}

func translatorConfigURL(translatorName, region string) (string, error) {
	config, err := loadTranslatorConfig(translatorName)
	if err != nil {
		return "", err
	}

	var cfg map[string]any
	if err := json.Unmarshal(config, &cfg); err != nil {
		return "", err
	}

	urls, ok := cfg["url"].(map[string]any)
	if !ok {
		return "", errors.New("url config is invalid")
	}
	url, ok := urls[region].(string)
	if !ok || strings.TrimSpace(url) == "" {
		return "", errors.New("url is invalid")
	}
	return url, nil
}
