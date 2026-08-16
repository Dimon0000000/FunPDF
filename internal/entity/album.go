package entity

type Album struct {
	ID          string `gorm:"primary_key" json:"id"`
	Name        string `json:"name"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
	BaseModel
}
