package entity

type AlbumFile struct {
	AlbumID string `json:"album_id"`
	FileID  string `json:"file_id"`
	//Deleted gorm.DeletedAt `gorm:"column:deleted_at"` TODO: do this later
	BaseModel
}
