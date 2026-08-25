package engine

import (
	"FunPDF/internal/common"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	BaiduTranslateURL           = "https://fanyi-api.baidu.com/ait/api/aiTextTranslate"
	BaiduTranslatorAPIReference = "https://fanyi-api.baidu.com/doc/21"
)

type BaiduTranslator struct {
	APIKey     string `json:"api_key"`
	APPID      string `json:"app_id"`
	httpClient *http.Client
}

func NewBaiduTranslator(apiKey string, appID string) *BaiduTranslator {
	return &BaiduTranslator{
		APIKey: apiKey,
		APPID:  appID,
		httpClient: &http.Client{
			Timeout: time.Second * 15,
		},
	}
}

// Translate use baidu translator translate text to dst language
func (b *BaiduTranslator) Translate(ctx context.Context, from, to, q string, params json.RawMessage) (string, error) {
	// build reqBody
	reqBody := make(map[string]any)

	q = strings.TrimSpace(q)
	if q == "" {
		return "", errors.New("the text is empty")
	}
	if len(q) > 6000 {
		return "", errors.New("the text is too long, it needs to smaller than 6000")
	}
	reqBody["q"] = ConcatenatingStrings(q)

	if from != "" {
		reqBody["from"] = from
	} else {
		reqBody["from"] = "auto"
	}

	if to != "" {
		reqBody["to"] = to
	} else {
		return "", errors.New("the dst language is empty")
	}

	reqBody["appid"] = b.APPID

	var external map[string]any
	if err := json.Unmarshal(params, &external); err != nil {
		return "", err
	}
	if external["model_type"] != nil && external["model_type"].(string) != "" {
		mType := external["model_type"].(string)
		if mType == "llm" || mType == "nmt" {
			reqBody["model_type"] = mType
		}
	}
	if external["reference"] != nil && external["reference"].(string) != "" {
		reqBody["reference"] = external["reference"].(string)
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// do request and get response
	toCtx, cancel := context.WithTimeout(ctx, b.httpClient.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(toCtx, "POST", BaiduTranslateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", b.APIKey))

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return "", errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("unmarshal response body: %w", err)
	}

	transResult, ok := result["trans_result"].([]any)
	if !ok {
		return "", fmt.Errorf("errors occurred: %s", result["error_msg"].(string))
	}

	var dst string
	for _, v := range transResult {
		dst = dst + v.(map[string]any)["dst"].(string)
	}

	return dst, nil
}

// Healthy check baidu translator is healthy
func (b *BaiduTranslator) Healthy(ctx context.Context) bool {
	_, err := b.Translate(ctx, "en", "zh", "hi", nil)
	common.Info(fmt.Sprintf("please read %s for more info", BaiduTranslatorAPIReference))
	return err == nil
}

// Name Baidu-Translator
func (b *BaiduTranslator) Name(ctx context.Context) string {
	return "Baidu-Translator"
}
