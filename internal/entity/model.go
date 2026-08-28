package entity

type Model struct {
	ID   string `gorm:"column:id;primaryKey;size:36" json:"id"`
	Name string `json:"name"`
	BaseModel
}

func (Model) TableName() string {
	return "models"
}
