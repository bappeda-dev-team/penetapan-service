package web

type ErrorResponse struct {
	Error any `json:"error"`
}

type ValidationErrorResponse struct {
	Error map[string]string `json:"error" example:"kodeOpd:required,tahun:required"`
}

type UnprocessableEntityResponse struct {
	Error map[string]string `json:"error" example:"kode_opd:format kode tidak valid"`
}
