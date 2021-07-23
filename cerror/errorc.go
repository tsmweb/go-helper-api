/*
Package cerror provides customized errors that represent status for
HTTP responses and specific treatments.

 */
package cerror

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest))
	ErrNotFound = errors.New(http.StatusText(http.StatusNotFound))
	ErrInternalServer = errors.New(http.StatusText(http.StatusInternalServerError))
	ErrUnauthorized = errors.New(http.StatusText(http.StatusUnauthorized))
)

// ErrValidateModel error issued when validating a model's fields.
type ErrValidateModel struct {
	Msg string
}

// Error implements interface Error
func (e *ErrValidateModel) Error() string {
	return e.Msg
}

// ErrRecordAlreadyRegistered indicates that the record already exists in the
// data structure in question.
var ErrRecordAlreadyRegistered = errors.New("Record Already Registered")
