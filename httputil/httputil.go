/*
Package httputil provides http utility methods.
*/

package httputil

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"
)

// MimeType represents the mime type ("application/json", "text/plain", "image/jpeg", ...).
type MimeType int

const (
	// MimeApplicationJSON represents the MimeType "application/json".
	MimeApplicationJSON = iota

	// MimeApplicationPDF represents the MimeType "application/pdf".
	MimeApplicationPDF

	// MimeApplicationOctetStream represents the MimeType "application/octet-stream".
	MimeApplicationOctetStream

	// MimeImageJPEG represents the MimeType "image/jpeg".
	MimeImageJPEG

	// MimeImagePNG represents the MimeType "image/png".
	MimeImagePNG

	// MimeAudioMP3 represents the MimeType "audio/mpeg".
	MimeAudioMP3

	// MimeVideoMP4 represents the MimeType "video/mp4".
	MimeVideoMP4

	// MimeTextPlain represents the MimeType "text/plain".
	MimeTextPlain
)

var mimeTypeText = map[MimeType]string{
	MimeApplicationJSON:        "application/json",
	MimeApplicationPDF:         "application/pdf",
	MimeApplicationOctetStream: "application/octet-stream",
	MimeImageJPEG:              "image/jpeg",
	MimeImagePNG:               "image/png",
	MimeAudioMP3:               "audio/mpeg",
	MimeVideoMP4:               "video/mp4",
	MimeTextPlain:              "text/plain",
}

// String return the name of the mimeType.
func (m MimeType) String() string {
	return mimeTypeText[m]
}

// MimeTypeText return the name of the mimeType.
func MimeTypeText(code MimeType) string {
	return code.String()
}

func HasContentType(r *http.Request, mimetype MimeType) bool {
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

// Headers http header.
type Headers map[string]string

func RespondWithHeader(w http.ResponseWriter, status int, h Headers) {
	for k, v := range h {
		w.Header().Set(k, v)
	}

	w.WriteHeader(status)
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithHeader(w, status, Headers{"Content-Type": MimeTypeText(MimeApplicationJSON)})
	w.Write([]byte(fmt.Sprintf(`{"error_message": "%s"}`, message)))
}

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	RespondWithHeader(w, status, Headers{"Content-Type": MimeTypeText(MimeApplicationJSON)})
	jsonData, err := json.Marshal(data)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	} else {
		w.Write(jsonData)
	}
}
