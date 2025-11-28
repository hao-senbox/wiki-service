package media_gateway_dto

type GetVideoUrlResponse struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	WikiCode        string `json:"wiki_code"`
	VideoURL        string `json:"video_url"`
	ImagePreviewURL string `json:"image_preview_url"`
	CreatedAt       string `json:"created_at"`
}
