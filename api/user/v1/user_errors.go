package v1

import (
	"github.com/go-kratos/kratos/v2/errors"
)

// User error definitions
var (
	ErrorUserNotFound        = errors.NotFound("USER_NOT_FOUND", "user not found")
	ErrorUserAlreadyExists   = errors.Conflict("USER_ALREADY_EXISTS", "user already exists")
	ErrorEmailAlreadyExists  = errors.Conflict("EMAIL_ALREADY_EXISTS", "email already exists")
	ErrorInvalidCredentials  = errors.Unauthorized("INVALID_CREDENTIALS", "invalid credentials")
	ErrorUserInactive        = errors.Forbidden("USER_INACTIVE", "user is inactive")
	ErrorInvalidPassword     = errors.BadRequest("INVALID_PASSWORD", "invalid password")
	ErrorPasswordTooWeak     = errors.BadRequest("PASSWORD_TOO_WEAK", "password too weak")
	ErrorInvalidEmailFormat  = errors.BadRequest("INVALID_EMAIL_FORMAT", "invalid email format")
	ErrorUsernameTooShort    = errors.BadRequest("USERNAME_TOO_SHORT", "username too short")
	ErrorUsernameTooLong     = errors.BadRequest("USERNAME_TOO_LONG", "username too long")
	ErrorPermissionDenied    = errors.Forbidden("PERMISSION_DENIED", "permission denied")
)

// Error constructor functions
func ErrorUserNotFoundWithMsg(msg string) *errors.Error {
	return errors.NotFound("USER_NOT_FOUND", msg)
}

func ErrorUserAlreadyExistsWithMsg(msg string) *errors.Error {
	return errors.Conflict("USER_ALREADY_EXISTS", msg)
}

func ErrorEmailAlreadyExistsWithMsg(msg string) *errors.Error {
	return errors.Conflict("EMAIL_ALREADY_EXISTS", msg)
}

func ErrorInvalidCredentialsWithMsg(msg string) *errors.Error {
	return errors.Unauthorized("INVALID_CREDENTIALS", msg)
}

func ErrorUserInactiveWithMsg(msg string) *errors.Error {
	return errors.Forbidden("USER_INACTIVE", msg)
}