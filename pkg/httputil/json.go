package httputil

import (
	"encoding/json"
	"net/http"
)

const contentType = "application/json; charset=utf-8"

func EncodeJSON(w http.ResponseWriter, data interface{}, code int) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func DecodeJSON(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}
