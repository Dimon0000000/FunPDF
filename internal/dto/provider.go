package dto

import "encoding/json"

type ListProvidersResult struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	BaseURL   string          `json:"url"`
	URLSuffix json.RawMessage `json:"url_suffix"`
}

type CreateProviderRequest struct {
	ID        string
	Name      string            `json:"name"`
	BaseURL   string            `json:"base_url"`
	URLSuffix map[string]string `json:"url_suffix"`
	APIKey    string            `json:"api_key"`
}

type UpdateProviderRequest struct {
	BaseURL   string            `json:"base_url"`
	URLSuffix map[string]string `json:"url_suffix"`
	APIKey    string            `json:"api_key"`
}
