package entity

type Provider struct {
	ID   string `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
	BaseModel
}

func (Provider) TableName() string {
	return "provider"
}
