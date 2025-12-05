package request

type UpdateWikiRequest struct {
	Language  *int      `json:"language"`
	Code      *string   `json:"code"`
	Public    *int      `json:"public"`
	Title     *string   `json:"title"`
	ImageWiki *string   `json:"image_wiki"`
	Keywords  *string   `json:"keywords"`
	Level     *int      `json:"level"`
	Unit      *string   `json:"unit"`
	Elements  []Element `json:"elements"`
}

type PictureItem struct {
	Key   string  `json:"key"`
	Order int     `json:"order"`
	Title *string `json:"title,omitempty"`
}

type Element struct {
	Number      int          `json:"number"`
	Type        string       `json:"type"`
	Value       *string      `json:"value,omitempty"`
	PictureKeys []PictureItem `json:"picture_keys,omitempty"`
	VideoID     *string      `json:"video_id,omitempty"`
}
