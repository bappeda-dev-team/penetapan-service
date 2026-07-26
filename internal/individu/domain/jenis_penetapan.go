package domain

type JenisPenetapan string

const (
	JenisPenetapanPk            JenisPenetapan = "PK-INDIVIDU"
	JenisPenetapanRenjaIndividu JenisPenetapan = "RENJA-INDIVIDU"
)

func (j JenisPenetapan) IsValid() bool {

	switch j {
	case
		JenisPenetapanPk:
		return true
	case
		JenisPenetapanRenjaIndividu:
		return true
	}
	return false

}
