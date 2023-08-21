package apperrors

import "encoding/json"

var (
	ErrNotFound = NewAppError(nil, "not found", "", "404")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
	Code             string `json:"code"`
}

func (e *AppError) Error() string {
	return e.Message
}
func (e *AppError) Unwrap() error {
	return e.Err
}
func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developerMessage, code string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}
func SystemError(err error) *AppError {
	return NewAppError(err, "system error", err.Error(), "500")
}
func BadRequestError(message string) *AppError {
	return NewAppError(nil, message, "bad request", "400")
}
