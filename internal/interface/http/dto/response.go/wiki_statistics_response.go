package response

import "time"

type WikiStatisticsResponse struct {
	Languages  map[string]string   `json:"languages"` // "1": "yes", "2": "no"
	Code       string              `json:"code"`
	Elements   []ElementStatistics `json:"elements"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

type ElementStatistics struct {
	Number int    `json:"number"`
	Type   string `json:"type"`
	Check  string `json:"check"` // "yes" or "no"
}

type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type WikiTemplateResponse struct {
	ID             string            `json:"id"`
	OrganizationID string            `json:"organization_id"`
	Type           string            `json:"type"`
	Name           string            `json:"name"`
	Description    *string           `json:"description"`
	Elements       []ElementResponse `json:"elements"`
	Icon           string            `json:"icon"`
	CreatedBy      string            `json:"created_by"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}
