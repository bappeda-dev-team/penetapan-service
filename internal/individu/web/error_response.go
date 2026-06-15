package web

type ErrorResponse struct {
	Error any `json:"error"`
}

type ValidationErrorResponse struct {
	Error map[string]string `json:"error" example:"pegawai_id:required,kode_opd:required,tahun:required"`
}
