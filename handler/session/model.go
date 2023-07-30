package session

const (
	daySeconds = 60 * 60 * 24
	CookieKey  = "session"
)

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	Uuid   string `json:"uuid"`
	Expiry int64  `json:"expiry"`
}
