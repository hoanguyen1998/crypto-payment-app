package helpers

import (
	"net/http"
)

// var (
// 	ErrAppNotFound           = errors.New("not found app with this id")
// 	ErrMasterKeyNotFound     = errors.New("master key for this app and this payment method not found")
// 	ErrAppKeyNotFound        = errors.New("application key for this app and this payment method not found")
// 	ErrAppKeyAlreadyExist    = errors.New("application key for this app and this payment method already exist")
// 	ErrMasterKeyAlreadyExist = errors.New("master key for this app and this payment method already exist")
// 	ErrOrderNotFound         = errors.New("not found order with this id")
// 	ErrNoPaymentMethodForApp = errors.New("this app doesn't have any payment method")
// )

var (
	ErrAppNotFound           = NewNotFoundError("not found app with this id")
	ErrMasterKeyNotFound     = NewNotFoundError("master key for this app and this payment method not found")
	ErrAppKeyNotFound        = NewNotFoundError("application key for this app and this payment method not found")
	ErrAppKeyAlreadyExist    = NewBadRequestError("application key for this app and this payment method already exist")
	ErrMasterKeyAlreadyExist = NewBadRequestError("master key for this app and this payment method already exist")
	ErrOrderNotFound         = NewNotFoundError("not found order with this id")
	ErrNoPaymentMethodForApp = NewBadRequestError("this app doesn't have any payment method")
	ErrEmailAlreadyExist     = NewBadRequestError("email is already exist")
)

const (
	serverError     = "server-error"
	notFoundError   = "not-found"
	badRequestError = "bad-request"
)

type ErrRest struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewInternalServerError(message string) *ErrRest {
	return &ErrRest{
		Status:  http.StatusInternalServerError,
		Error:   serverError,
		Message: message,
	}
}

func NewBadRequestError(message string) *ErrRest {
	return &ErrRest{
		Status:  http.StatusBadRequest,
		Error:   badRequestError,
		Message: message,
	}
}

func NewNotFoundError(message string) *ErrRest {
	return &ErrRest{
		Status:  http.StatusNotFound,
		Error:   notFoundError,
		Message: message,
	}
}

func NewForbiddenError(message string) *ErrRest {
	return &ErrRest{
		Status:  http.StatusForbidden,
		Error:   serverError,
		Message: message,
	}
}
