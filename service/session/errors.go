package session

type UserNotVerifiedError struct {
}

func (u UserNotVerifiedError) Error() string {
	return "User not verified"
}

type LicenseExpiredError struct{}

func (l LicenseExpiredError) Error() string {
	return "License has expired"
}
