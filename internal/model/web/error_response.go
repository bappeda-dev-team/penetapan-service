package web

type ErrorResponse struct {
	Error any `json:"error"`
}

type ValidationErrorResponse struct {
	Error map[string]string `json:"error"`
}
