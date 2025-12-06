package request

type CreateWikiTemplateRequest struct {
	Type     string    `json:"type"`
	Elements []Element `json:"elements"`
	Status   string    `json:"status"`
}
