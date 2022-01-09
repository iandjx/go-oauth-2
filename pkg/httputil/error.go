package httputil

import (
	"net/http"
)

// Error400 reply to the request with string value of error and status code of bad request.
func Error400(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusBadRequest)
}

// Error403 reply to the request with string value of error and status code of forbidden.
func Error403(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusForbidden)
}

// Error404 reply to the request with string value of error and status code of not found.
func Error404(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusNotFound)
}

// Error500 reply to the request with string value of error and status code of internal error.
func Error500(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusInternalServerError)
}

// Error reply to the request with string value of error and status code.
func Error(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}
