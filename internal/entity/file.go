package enetity

type File struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	MimeType          string `json:"mime_type"`
	ProjectStorageKey string `json:"-"`
	Size              int64  `json:"size"`
	SHA256            string `json:"sha256"`
	Revision          int64  `json:"revision"`
	Status            string `json:"status"`
	BaseModel
}
