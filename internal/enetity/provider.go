package enetity

type Provider struct {
	ID   string
	Name string
	BaseModel
}

func (Provider) TableName() string {
	return "provider"
}
