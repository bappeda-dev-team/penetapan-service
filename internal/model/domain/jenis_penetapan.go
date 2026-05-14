package domain

type JenisPenetapan string

const (
	JenisPenetapanTujuan  = "TUJUAN-OPD"
	JenisPenetapanSasaran = "SASARAN-OPD"
)

func (j JenisPenetapan) IsValid() bool {

	switch j {
	case
		JenisPenetapanTujuan,
		JenisPenetapanSasaran:
		return true
	}
	return false

}
