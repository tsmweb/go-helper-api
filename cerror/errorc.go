/*
Package cerror provides customized errors that represent status for
HTTP responses and specific treatments.

 */
package cerror

import (
	"errors"
	"net/http"
)

var ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest))

var ErrNotFound = errors.New(http.StatusText(http.StatusNotFound))

var ErrInternalServer = errors.New(http.StatusText(http.StatusInternalServerError))

var ErrUnauthorized = errors.New(http.StatusText(http.StatusUnauthorized))

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
