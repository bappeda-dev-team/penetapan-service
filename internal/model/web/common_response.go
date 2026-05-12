package web

type TargetIndikatorResponse struct {
	Id       int64   `json:"id" example:"1"`
	Tahun    int     `json:"tahun" example:"2025"`
	Target   float64 `json:"target" example:"100"`
	Satuan   string  `json:"satuan" example:"%"`
}
