package response

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
