package domain

type JenisPenetapan string

const (
	JenisPenetapanTujuanPemda JenisPenetapan = "TUJUAN-PEMDA"
)

func (j JenisPenetapan) IsValid() bool {

	switch j {
	case
		JenisPenetapanTujuanPemda:
		return true
	}
	return false

}
