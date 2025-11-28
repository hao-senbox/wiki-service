package media_gateway_dto

type GetVideoUrlRequest struct {
	VideoID  string `json:"video_id" binding:"required"`
	Language *int   `json:"language" binding:"required"`
}
