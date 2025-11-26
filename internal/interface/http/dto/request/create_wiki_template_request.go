package request

type CreateWikiTemplateRequest struct {
	OrganizationID string    `json:"organization_id"`
	Type           string    `json:"type"`
	Elements       []Element `json:"elements"`
}

type Element struct {
	Number  int     `json:"number"`
	Type    string  `json:"type"`
	Value   *string `json:"value"`
	TopicID *string `json:"topic_id,omitempty"`
}
