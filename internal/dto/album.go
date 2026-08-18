package dto

type CreateAlbumRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
}

type UpdateAlbumRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
}

type AlertAlbumFilesRequest struct {
	IDs []string `json:"ids"`
}
