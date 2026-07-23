package common

type ErrorType string

const (
	Validation ErrorType = "validation"
	NotFound   ErrorType = "not_found"
	Conflict   ErrorType = "conflict"
	Internal   ErrorType = "internal"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e AppError) Unwrap() error {
	return e.Err
}

func NewValidation(msg string) error {
	return AppError{
		Type:    Validation,
		Message: msg,
	}
}

func WrapInternal(msg string, err error) error {
	return AppError{
		Type:    Internal,
		Message: msg,
		Err:     err,
	}
}
