package dto

type ListModelsResponse struct {
	Name string `json:"name"`
}

type ChatResponse struct {
	Answer        *string `json:"answer"`
	ReasonContent *string `json:"reason_content"`
}

type ChatRequest struct {
	ModelName   string           `json:"model_name"`
	ModelID     string           `json:"model_id"`
	Messages    []map[string]any `json:"messages"`
	Stream      *bool            `json:"stream"`
	Temperature *float64         `json:"temperature"`
	TopP        *float64         `json:"top_p"`
	Thinking    *bool            `json:"thinking"`
	MaxTokens   *int             `json:"max_tokens"`
	Effort      *string          `json:"effort"`
	Verbosity   *string          `json:"verbosity"`
}
