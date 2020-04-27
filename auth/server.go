package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

func New(prefix string) http.Handler {
	t := BuiltinTemplates()
	m := mux.NewRouter()
	l := &Login{
		t: t,
	}
	l.registerRoutes(m.PathPrefix(prefix).Subrouter())
	return m
}
