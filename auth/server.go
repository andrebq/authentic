package auth

import (
	"net/http"

	"github.com/andrebq/authentic/internal/session"
	"github.com/gorilla/mux"
)

func New(prefix string, s *session.S, cat UserCatalog) http.Handler {
	t := BuiltinTemplates()
	m := mux.NewRouter()
	l := &Login{
		t:       t,
		s:       s,
		catalog: cat,
	}
	l.registerRoutes(m.PathPrefix(prefix).Subrouter())
	return m
}
