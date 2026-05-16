package domain

type JenisPenetapan string

const (
	JenisPenetapanRenaksi = "RENAKSI-OPD"
	JenisPenetapanRenja   = "RENJA-OPD"
	JenisPenetapanTujuan  = "TUJUAN-OPD"
	JenisPenetapanSasaran = "SASARAN-OPD"
)

func (j JenisPenetapan) IsValid() bool {

	switch j {
	case
		JenisPenetapanRenaksi,
		JenisPenetapanRenja,
		JenisPenetapanTujuan,
		JenisPenetapanSasaran:
		return true
	}
	return false

}
