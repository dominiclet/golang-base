package user

type SendVerificationEmailErrType int

const (
	ResendDisabled SendVerificationEmailErrType = iota
	UserVerified
)

type SendVerificationEmailError struct {
	msg     string
	errType SendVerificationEmailErrType
}

func NewSendVerificationEmailErr(errType SendVerificationEmailErrType, msg string) error {
	return SendVerificationEmailError{
		msg:     msg,
		errType: errType,
	}
}

func (s SendVerificationEmailError) Error() string {
	return s.msg
}

func (s SendVerificationEmailError) ErrType() SendVerificationEmailErrType {
	return s.errType
}
