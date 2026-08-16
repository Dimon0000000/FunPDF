package entity

type Model struct {
	ID           string `gorm:"column:id;primaryKey;size:32" json:"id"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	ApiKey       string `json:"-"`
	BaseModel
}

func (Model) TableName() string {
	return "models"
}
