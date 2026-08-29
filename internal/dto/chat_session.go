package dto

type SetupChatSessionRequest struct {
	FileID       string `json:"file_id"`
	ModelID      string `json:"model_id"`
	ModelName    string `json:"model_name"`
	SystemPrompt string `json:"system_prompt"`
}

type UpdateDialogStatusRequest struct {
	Status int `json:"status"`
}

type SendMessageRequest struct {
	Content     string   `json:"content"`
	Quote       string   `json:"quote"`
	Stream      *bool    `json:"stream"`
	Temperature *float64 `json:"temperature"`
	TopP        *float64 `json:"top_p"`
	Thinking    *bool    `json:"thinking"`
	MaxTokens   *int     `json:"max_tokens"`
	Effort      *string  `json:"effort"`
}
