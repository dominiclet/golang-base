package resperror

import "net/http"

var ErrorMapping map[int]CustomErrWithCode = map[int]CustomErrWithCode{
	// General errors
	BadRequest: {
		StatusCode: http.StatusBadRequest,
		Code:       BadRequest,
		Message:    "Invalid request",
	},
	TooManyRequests: {
		StatusCode: http.StatusTooManyRequests,
		Code:       TooManyRequests,
		Message:    "Too many requests sent, try again later",
	},
	Unauthorized: {
		StatusCode: http.StatusUnauthorized,
		Code:       Unauthorized,
		Message:    "Unauthorized",
	},
	// User errors
	UserAlreadyExistsError: {
		StatusCode: http.StatusConflict,
		Code:       UserAlreadyExistsError,
		Message:    "User with matching email already exists",
	},
	UserNotFound: {
		StatusCode: http.StatusNotFound,
		Code:       UserNotFound,
		Message:    "User not found",
	},
	UserNotVerifiedError: {
		StatusCode: http.StatusForbidden,
		Code:       UserNotVerifiedError,
		Message:    "Account email is not verified",
	},
	UserLicenseExpiredError: {
		StatusCode: http.StatusMethodNotAllowed,
		Code:       UserLicenseExpiredError,
		Message:    "License expired",
	},
	UserEmailNotFound: {
		StatusCode: http.StatusNotFound,
		Code:       UserEmailNotFound,
		Message:    "An account with the provided email cannot be found",
	},
	UserIncorrectPassword: {
		StatusCode: http.StatusUnauthorized,
		Code:       UserIncorrectPassword,
		Message:    "Incorrect password",
	},
	// Email verification errors
	UserAlreadyVerifiedError: {
		StatusCode: http.StatusMethodNotAllowed,
		Code:       UserAlreadyVerifiedError,
		Message:    "Failed to send verification email: user already verified",
	},
	InvalidTokenError: {
		StatusCode: http.StatusUnauthorized,
		Code:       InvalidTokenError,
		Message:    "Invalid token",
	},
}
