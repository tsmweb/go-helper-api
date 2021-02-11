/*
Package controller implements a Controller with utility methods.
The Controller is used in the composition of other more specific controllers.

Package controller also provides the MimeType type which represents the mime type as
"application/json", "text/plain", "image/jpeg", ...

 */
package controller

import (
	"encoding/json"
	"github.com/tsmweb/go-helper-api/auth"
	"mime"
	"net/http"
	"strings"
)

type errorMessage struct {
	ErrorMessage string `json:"error_message"`
}

type Headers map[string]string

// Controller base controller with utility methods.
type Controller struct {
	jwt auth.JWT
}

// New returns an instance of the Controller.
func New(jwt auth.JWT) *Controller {
	return &Controller{jwt}
}

// ExtractID extracts the JWT token id.
func (c *Controller) ExtractID(r *http.Request) (string, error) {
	ID, err := c.jwt.GetDataToken(r, "id")
	if err != nil || ID == nil {
		return "", err
	}

	return ID.(string), nil
}

// HasContentType validates the content type, return true for valid and false for invalid.
func (c *Controller) HasContentType(r *http.Request, mimetype MimeType) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return mimetype == MimeApplicationOctetStream
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}

		if t == mimetype.String() {
			return true
		}
	}

	return false
}

func (c *Controller) RespondWithHeader(w http.ResponseWriter, status int, h Headers) {
	for k, v := range h {
		w.Header().Set(k, v)
	}

	w.WriteHeader(status)
}

func (c *Controller) RespondWithError(w http.ResponseWriter, status int, message string) {
	c.RespondWithHeader(w, status, Headers{"Content-Type": MimeTypeText(MimeApplicationJSON)})
	errorMessage, err := json.Marshal(errorMessage{ErrorMessage: message})
	if err == nil {
		w.Write(errorMessage)
	}
}

func (c *Controller) RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	c.RespondWithHeader(w, status, Headers{"Content-Type": MimeTypeText(MimeApplicationJSON)})
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	} else {
		w.Write(jsonData)
	}
}
