package helper

import (
	"strings"
	"unicode"
)

func TextCleaner(text string) string {
	text = strings.Map(func(r rune) rune {
		switch {
		// Hapus karakter kontrol
		case unicode.IsControl(r):
			return -1

		// Hilangkan semua jenis quotation mark Unicode
		case unicode.In(r, unicode.Quotation_Mark):
			return -1

		// Samakan semua whitespace menjadi spasi
		case unicode.IsSpace(r):
			return ' '

		default:
			return r
		}
	}, text)

	// Trim spasi depan, belakang, dan rapikan spasi di tengah
	return strings.Join(strings.Fields(text), " ")
}
