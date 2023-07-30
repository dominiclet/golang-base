package resperror

import "net/http"

type CustomErrWithCode struct {
	Code       int
	Message    string
	StatusCode int
}

func (c CustomErrWithCode) Error() string {
	return c.Message
}

func NewError(code int) CustomErrWithCode {
	customErr, ok := ErrorMapping[code]
	if !ok {
		return CustomErrWithCode{
			Code:       code,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return customErr
}
