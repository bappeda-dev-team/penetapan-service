package helper

import (
	"regexp"
	"strconv"
	"strings"
)

var numberPattern = regexp.MustCompile(`[-+]?\d+(?:[.,]\d+)?`)

func ParseFloat(value string) (float64, error) {
	value = strings.TrimSpace(value)

	// Ambil angka pertama dari string.
	// Contoh:
	// "61,10 / B" -> "61,10"
	// "61.10 / B" -> "61.10"
	// "61 / B"    -> "61"
	match := numberPattern.FindString(value)

	if match == "" {
		return 0, strconv.ErrSyntax
	}

	// Standarisasi decimal separator Indonesia -> Go float.
	match = strings.Replace(match, ",", ".", 1)

	return strconv.ParseFloat(match, 64)
}

func ParseTahun(value string) (int, error) {
	value = strings.TrimSpace(value)

	match := numberPattern.FindString(value)

	if match == "" {
		return 0, strconv.ErrSyntax
	}

	return strconv.Atoi(value)
}
