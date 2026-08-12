package helper

import (
	"regexp"
	"strconv"
	"strings"
)

var numberPattern = regexp.MustCompile(`[-+]?\d+(?:[.,]\d+)?`)

func ParseTahun(value string) (int, error) {
	value = strings.TrimSpace(value)

	match := numberPattern.FindString(value)

	if match == "" {
		return 0, strconv.ErrSyntax
	}

	return strconv.Atoi(value)
}
