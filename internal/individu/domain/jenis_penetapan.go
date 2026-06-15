package domain

type JenisPenetapan string

const (
	JenisPenetapanPk JenisPenetapan = "PK-INDIVIDU"
)

func (j JenisPenetapan) IsValid() bool {

	switch j {
	case
		JenisPenetapanPk:
		return true
	}
	return false

}
