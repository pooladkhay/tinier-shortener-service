package errs

import (
	"errors"
	"net/http"
)

type Err struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewErr(msg string) error {
	return errors.New(msg)
}

func NewBadRequestError(msg string) *Err {
	return &Err{
		Message: msg,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}
func NewUnauthorizedError(msg string) *Err {
	return &Err{
		Message: msg,
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}

func NewNotFoundError(msg string) *Err {
	return &Err{
		Message: msg,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

func NewConflictError(msg string) *Err {
	return &Err{
		Message: msg,
		Status:  http.StatusConflict,
		Error:   "conflict",
	}
}

func NewInternalServerError(msg string) *Err {
	return &Err{
		Message: msg,
		Status:  http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
}
