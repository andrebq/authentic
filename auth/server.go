package auth

import (
	"net/http"

	"github.com/andrebq/authentic/internal/session"
	"github.com/gorilla/mux"
)

func New(prefix string, s *session.S) http.Handler {
	t := BuiltinTemplates()
	m := mux.NewRouter()
	l := &Login{
		t: t,
		s: s,
	}
	l.registerRoutes(m.PathPrefix(prefix).Subrouter())
	return m
}
