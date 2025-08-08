package apperrors

type ServerError struct {
	Code    int    // HTTP-статус код
	Message string // Сообщение для клиента
	Err     error  // Оригинальная ошибка (для логирования)
}

func (e *ServerError) Error() string {
	return e.Message
}

func New(code int, message string, err error) *ServerError {
	return &ServerError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrInvalidInput  = New(400, "Invalid input data", nil)
	ErrInternal      = New(500, "Internal server error", nil)
	ErrAlreadyExists = New(503, "event already exists", nil)
	ErrNotFound      = New(503, "Event not found", nil)
	ErrPastDate      = New(503, "Date cannot be in the past", nil)
)
