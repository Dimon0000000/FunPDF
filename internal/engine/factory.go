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

func (t *TranslatorFactory) GetTranslator(ctx context.Context, db *gorm.DB, translatorName string) (Translator, error) {
	translatorName = strings.TrimSpace(translatorName)
	if translatorName == "" {
		return nil, errors.New("translator name is required")
	}

	params, err := t.translatorDAO.GetTranslatorParams(ctx, db, translatorName)
	if err != nil {
		return nil, err
	}
	var param map[string]any
	err = json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}

	switch translatorName {
	case "Baidu-Translator":
		apiKey, ok := param["api_key"].(string)
		if !ok {
			return nil, errors.New("api_key is invalid")
		}
		appID, ok := param["app_id"].(string)
		if !ok {
			return nil, errors.New("app_id is invalid")
		}
		return NewBaiduTranslator(apiKey, appID), nil
	default:
		return nil, fmt.Errorf("unsupported translator: %s", translatorName)
	}
}
