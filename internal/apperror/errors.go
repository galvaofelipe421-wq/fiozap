package apperror

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Cause      error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) WithCause(err error) *Error {
	return &Error{
		Code:       e.Code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Cause:      err,
	}
}

func (e *Error) WithMessage(msg string) *Error {
	return &Error{
		Code:       e.Code,
		Message:    msg,
		StatusCode: e.StatusCode,
		Cause:      e.Cause,
	}
}

// Authentication errors
var (
	ErrMissingToken = &Error{
		Code:       "MISSING_TOKEN",
		Message:    "missing authentication token",
		StatusCode: http.StatusUnauthorized,
	}
	ErrInvalidToken = &Error{
		Code:       "INVALID_TOKEN",
		Message:    "invalid authentication token",
		StatusCode: http.StatusUnauthorized,
	}
	ErrUnauthorized = &Error{
		Code:       "UNAUTHORIZED",
		Message:    "unauthorized",
		StatusCode: http.StatusUnauthorized,
	}
)

// Session errors
var (
	ErrSessionNotFound = &Error{
		Code:       "SESSION_NOT_FOUND",
		Message:    "session not found",
		StatusCode: http.StatusNotFound,
	}
	ErrSessionRequired = &Error{
		Code:       "SESSION_REQUIRED",
		Message:    "session name is required",
		StatusCode: http.StatusBadRequest,
	}
	ErrSessionLimitReached = &Error{
		Code:       "SESSION_LIMIT_REACHED",
		Message:    "session limit reached",
		StatusCode: http.StatusForbidden,
	}
	ErrAlreadyConnected = &Error{
		Code:       "ALREADY_CONNECTED",
		Message:    "session already connected",
		StatusCode: http.StatusConflict,
	}
	ErrNotConnected = &Error{
		Code:       "NOT_CONNECTED",
		Message:    "session not connected",
		StatusCode: http.StatusBadRequest,
	}
)

// User errors
var (
	ErrUserNotFound = &Error{
		Code:       "USER_NOT_FOUND",
		Message:    "user not found",
		StatusCode: http.StatusNotFound,
	}
)

// Validation errors
var (
	ErrInvalidPayload = &Error{
		Code:       "INVALID_PAYLOAD",
		Message:    "invalid request payload",
		StatusCode: http.StatusBadRequest,
	}
	ErrPhoneRequired = &Error{
		Code:       "PHONE_REQUIRED",
		Message:    "phone is required",
		StatusCode: http.StatusBadRequest,
	}
	ErrMessageRequired = &Error{
		Code:       "MESSAGE_REQUIRED",
		Message:    "message is required",
		StatusCode: http.StatusBadRequest,
	}
	ErrInvalidPhone = &Error{
		Code:       "INVALID_PHONE",
		Message:    "invalid phone number format",
		StatusCode: http.StatusBadRequest,
	}
)

// WhatsApp errors
var (
	ErrNoSession = &Error{
		Code:       "NO_WHATSAPP_SESSION",
		Message:    "no active WhatsApp session",
		StatusCode: http.StatusServiceUnavailable,
	}
	ErrSendFailed = &Error{
		Code:       "SEND_FAILED",
		Message:    "failed to send message",
		StatusCode: http.StatusInternalServerError,
	}
	ErrUploadFailed = &Error{
		Code:       "UPLOAD_FAILED",
		Message:    "failed to upload media",
		StatusCode: http.StatusInternalServerError,
	}
)

// Internal errors
var (
	ErrInternal = &Error{
		Code:       "INTERNAL_ERROR",
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
	}
	ErrDatabaseError = &Error{
		Code:       "DATABASE_ERROR",
		Message:    "database error",
		StatusCode: http.StatusInternalServerError,
	}
)

func New(code, message string, statusCode int) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func BadRequest(message string) *Error {
	return &Error{
		Code:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NotFound(message string) *Error {
	return &Error{
		Code:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func Internal(message string) *Error {
	return &Error{
		Code:       "INTERNAL_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}
