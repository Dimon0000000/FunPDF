package entity

type AlbumFile struct {
	AlbumID   string `json:"album_id"`
	FileID    string `json:"file_id"`
	SortOrder int    `json:"sort_order"`
	BaseModel
}
