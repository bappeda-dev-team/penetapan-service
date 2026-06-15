package validator

import (
	"regexp"
	"slices"
)

var (
	KodeOpdRX   = regexp.MustCompile(`^\d\.\d{2}\.\d\.\d{2}\.\d\.\d{2}\.\d{2}\.\d{4}$`)
	PegawaiIdRX = regexp.MustCompile(`^\d{18}(_plt)?$`)
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, msg string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = msg
	}
}

// Check is used to check if the request is valid.
// pakai value yang valid untuk dicek
func (v *Validator) Check(ok bool, key, msg string) {
	if !ok {
		v.AddError(key, msg)
	}
}

// PermittedValue is used as enums for request
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Generic function which returns true if all values in a slice are unique.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
