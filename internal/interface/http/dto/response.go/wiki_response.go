package response

import "time"

type WikiResponse struct {
	ID             string                `json:"id"`
	OrganizationID string                `json:"organization_id"`
	Code           string                `json:"code"`
	Translation    []TranslationResponse `json:"translation"`
	Icon           string                `json:"icon"`
	CreatedBy      string                `json:"created_by"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
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
