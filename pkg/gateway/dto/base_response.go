package response

// APIGateWayResponse generic response structure from gateway
type APIGateWayResponse[T any] struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}
