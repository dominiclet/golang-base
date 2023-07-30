package resperror

// Error codes are in the format XYYZZ, where X represents the error code version,
// YY represents the module where the error happened, and ZZ enumerates the error which occurred

// General errors
const (
	UnknownError    = 10000
	BadRequest      = 10001
	TooManyRequests = 10002 // Rate limit
	Unauthorized    = 10003
)

// User
const (
	UserAlreadyExistsError  = 10101
	UserNotFound            = 10102
	UserNotVerifiedError    = 10103
	UserLicenseExpiredError = 10104
	UserEmailNotFound       = 10105
	UserIncorrectPassword   = 10106
)

// Email verification
const (
	UserAlreadyVerifiedError = 10201
	InvalidTokenError        = 10202
)
