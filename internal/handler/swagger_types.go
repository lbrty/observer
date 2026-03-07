package handler

// ErrorResponse is a generic error response.
type ErrorResponse struct {
	Error string `json:"error" example:"resource not found"`
	Code  string `json:"code" example:"errors.internal"`
}

// MessageResponse is a generic success message.
type MessageResponse struct {
	Message string `json:"message" example:"deleted"`
}
