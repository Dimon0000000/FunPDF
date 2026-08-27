package dto

type ListProvidersResult struct {
	Name    string            `json:"name"`
	BaseURL map[string]string `json:"url"`
}

type CreateProviderRequest struct {
	Name      string            `json:"name"`
	BaseURL   string            `json:"base_url"`
	URLSuffix map[string]string `json:"url_suffix"`
	APIKey    string            `json:"api_key"`
}

type UpdateProviderRequest struct {
	BaseURL   string `json:"base_url"`
	URLSuffix string `json:"url_suffix"`
	APIKey    string `json:"api_key"`
}
