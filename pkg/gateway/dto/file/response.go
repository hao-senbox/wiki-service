package file_gateway_dto

type UploadAudioResponse struct {
	Key string `json:"key"`
}

type UploadImageResponse struct {
	ImageName string `json:"image_name"`
	Key       string `json:"key"`
	Extension string `json:"extension"`
	Url       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type UploadPDFResponse struct {
	Key string `json:"key"`
}

type UploadVideoResponse struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type Avatar struct {
	ImageID  uint64 `json:"image_id"`
	ImageKey string `json:"image_key"`
	ImageUrl string `json:"image_url"`
	Index    int    `json:"index"`
	IsMain   bool   `json:"is_main"`
}
