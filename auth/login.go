package auth

import (
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/andrebq/authentic/internal/session"
	"github.com/andrebq/authentic/internal/webflow"
	"github.com/gorilla/mux"
)

type (
	Login struct {
		once    sync.Once
		t       *template.Template
		s       *session.S
		catalog UserCatalog
		realm   string
	}

	// UserCatalog implements a read-only database of user information
	UserCatalog interface {
		// Authenticate re
		Authenticate(username, password string) error
	}
)

// New login page - POST
func (l *Login) New(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = l.catalog.Authenticate(req.Form.Get("username"), req.Form.Get("password"))
	if err != nil {
		webflow.Authenticate(w, req, l.realm, "Invalid credentials")
		return
	}

	tk, expire, err := l.s.Start(time.Now())
	if err != nil {
		webflow.InternalError(w, req)
		return
	}
	cookie := &http.Cookie{}
	cookie.Name = "_session"
	cookie.Value = tk
	cookie.Expires = expire
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Domain = req.URL.Host
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	http.SetCookie(w, cookie)
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
	m.Methods("POST").HandlerFunc(l.New).Name("login_new")
}
