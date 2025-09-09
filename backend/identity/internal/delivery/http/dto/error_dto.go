package dto

type ErrorResponse struct {
	StatusCode int         `json:"status_code"`
	ErrorType  string      `json:"error_type"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}
