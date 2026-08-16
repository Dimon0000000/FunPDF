package entity

type File struct {
	ID             string `gorm:"primary_key" json:"id"`
	Name           string `json:"name"`
	MimeType       string `json:"mime_type"`
	FileStorageKey string `json:"-"`
	Size           int64  `json:"size"`
	SHA256         string `json:"sha256"`
	Revision       int64  `json:"revision"`
	Status         string `json:"status"`
	BaseModel
}
