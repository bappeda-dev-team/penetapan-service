package web

type TargetIndikatorResponse struct {
	Id       int64   `json:"id"`
	Tahun    int     `json:"tahun"`
	Target   float64 `json:"target"`
	Satuan   string  `json:"satuan"`
}
