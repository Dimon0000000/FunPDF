package entity

type ChatSession struct {
	ID           string `gorm:"primary_key" json:"id"`
	FileID       string `gorm:"index" json:"file_id"`
	ProviderID   string `gorm:"index" json:"provider_id"`
	ModelID      string `gorm:"index" json:"model_id"`
	ModelName    string `json:"model_name"`
	SystemPrompt string `gorm:"type:text" json:"-"`
	BaseModel
}
