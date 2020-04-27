package auth

import (
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type (
	Login struct {
		once sync.Once
		t    *template.Template
	}
)

// New login page - POST
func (l *Login) New(w http.ResponseWriter, req *http.Request) {
	cookie := http.Cookie{}
	cookie.Name = "_session"
	cookie.Value = "hello"
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Domain = req.URL.Host
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Expires = time.Now().Add(time.Hour)
	w.Header().Add("Set-Cookie", cookie.String())
	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusSeeOther)
}

// Create a new login - GET
func (l *Login) Create(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	l.t.ExecuteTemplate(w, "login/new", l)
}

func (l *Login) registerRoutes(m *mux.Router) {
	m = m.PathPrefix("/login").Subrouter()
	m.Methods("GET").HandlerFunc(l.Create).Name("login_create")
	m.Methods("POST").HandlerFunc(l.New).Name("login_create")
}
