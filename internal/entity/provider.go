package entity

type Provider struct {
	ID        string            `gorm:"primary_key" json:"id"`
	Name      string            `json:"name"`
	BaseURL   string            `json:"base_url"`
	URLSuffix map[string]string `gorm:"serializer:json" json:"url_suffix"`
	APIKey    string            `json:"api_key"`
	BaseModel
}

func (Provider) TableName() string {
	return "provider"
}
