package webflow

import (
	"io"
	"net/http"
)

// Authenticate responds with a 401 and configures the required headers
func Authenticate(w http.ResponseWriter, req *http.Request, realm, location string) {
	w.Header().Add("WWW-Authenticate", realm)
	w.WriteHeader(http.StatusUnauthorized)
	io.WriteString(w, location)
}

// InternalError responds with a 500 error
func InternalError(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "unexpected error")
}
