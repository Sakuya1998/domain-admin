package errors

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string) error {
	return &AppError{Code: code, Message: msg}
}
