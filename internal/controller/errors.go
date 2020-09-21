package controller

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Code int64
	err  error
}

func NewError(code int64, err error) *Error {
	return &Error{
		Code: code,
		err:  err,
	}
}

////////////////////////////////////////////////////////////////////////////////

func (E *Error) Error() string {
	if E.err == nil {
		return ""
	}

	return fmt.Sprintf("%d: %s", E.Code, E.err.Error())
}

var (
	ErrAccountNotFound = NewError(http.StatusNotFound, errors.New("account not found"))
	ErrBalanceNotOk    = NewError(http.StatusUnavailableForLegalReasons, errors.New("balance not ok"))
	ErrWrongPassword   = NewError(http.StatusNotAcceptable, errors.New("wrong password"))
	ErrInternal        = NewError(http.StatusInternalServerError, errors.New("server error"))
	ErrInvalidParam    = NewError(http.StatusBadRequest, errors.New("invalid param"))
	ErrUnauthorized    = NewError(http.StatusUnauthorized, errors.New("unauthorize or authorize-fail"))
)
