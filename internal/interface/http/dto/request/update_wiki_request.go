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

type Element struct {
	Number      int      `json:"number"`
	Type        string   `json:"type"`
	Value       *string  `json:"value,omitempty"`
	PictureKeys []string `json:"picture_keys,omitempty"`
	TopicID     *string  `json:"topic_id,omitempty"`
}
