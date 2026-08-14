package enetity

type Provider struct {
	ID   string
	Name string
	Base
}

func (Provider) TableName() string {
	return "provider"
}
