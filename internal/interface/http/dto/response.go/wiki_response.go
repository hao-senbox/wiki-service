package response

import "time"

type WikiResponse struct {
	ID            string                `json:"id"`
	Code          string                `json:"code"`
	Public        int                   `json:"public"`
	Translation   []TranslationResponse `json:"translation"`
	ImageWiki     string                `json:"image_wiki"`
	CreatedByUser *CreatedByUserInfo    `json:"creator,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

type CreatedByUserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type TranslationResponse struct {
	Language *int              `json:"language"`
	Title    *string           `json:"title"`
	Keywords *string           `json:"keywords"`
	Level    *int              `json:"level"`
	Unit     *string           `json:"unit"`
	Elements []ElementResponse `json:"elements"`
}

type ElementResponse struct {
	Number         int                `json:"number"`
	Type           string             `json:"type"`
	Value          *string            `json:"value,omitempty"`
	ImageUrl       *string            `json:"image_url,omitempty"`
	PdfUrl         *string            `json:"pdf_url,omitempty"`
	PictureKeys    []string           `json:"picture_keys,omitempty"`
	PictureKeysUrl []string           `json:"picture_keys_url,omitempty"`
	ButtonUrl      *ButtonUrlResponse `json:"button_url,omitempty"`
	Button         *ButtonResponse    `json:"button,omitempty"`
	VideoID        *string            `json:"video_id,omitempty"`
}

type ButtonUrlResponse struct {
	Title         string `json:"title"`
	ButtonUrl     string `json:"button_url"`
	ButtonIcon    string `json:"button_icon"`
	ButtonIconUrl string `json:"button_icon_url"`
}

type ButtonResponse struct {
	Title         string `json:"title"`
	ButtonIcon    string `json:"button_icon"`
	ButtonIconUrl string `json:"button_icon_url"`
}
