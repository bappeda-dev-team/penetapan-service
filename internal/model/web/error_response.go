package web

type ErrorResponse struct {
	Error any `json:"error"`
}

type ValidationErrorResponse struct {
	Error map[string]string `json:"error" example:"kode_opd:required,tahun:required"`
}
