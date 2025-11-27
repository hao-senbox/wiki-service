package response

import "time"

type WikiResponse struct {
	ID            string                `json:"id"`
	Code          string                `json:"code"`
	Public        int                   `json:"public"`
	Translation   []TranslationResponse `json:"translation"`
	Icon          string                `json:"icon"`
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
	Number  int     `json:"number"`
	Type    string  `json:"type"`
	Value   *string `json:"value"`
	TopicID *string `json:"topic_id"`
}
