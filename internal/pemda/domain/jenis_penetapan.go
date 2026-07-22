package domain

type JenisPenetapan string

const (
	JenisPenetapanTujuanPemda  JenisPenetapan = "TUJUAN-PEMDA"
	JenisPenetapanSasaranPemda JenisPenetapan = "SASARAN-PEMDA"
)

func (j JenisPenetapan) IsValid() bool {

	switch j {
	case
		JenisPenetapanTujuanPemda,
		JenisPenetapanSasaranPemda:
		return true
	}
	return false

}
