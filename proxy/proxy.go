package proxy

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type (
	// Reverse proxy with cookie checks
	Reverse struct {
		actual     *httputil.ReverseProxy
		cookieName string
		realm      string
		tokens     TokenSet
	}

	// TokenSet represents the entire set of valid tokens
	TokenSet interface {
		// Contains returns true if the token is still valid
		Contains(string) (bool, error)
	}
)

//NewReverse proxy which checks a specific cookie for protection,
// if the cookie is not present any request will return 401.
//
// If the cookie is present, access to the target is allowed
func NewReverse(cookieName, realm string, tokens TokenSet, target *url.URL) *Reverse {
	return &Reverse{
		actual:     httputil.NewSingleHostReverseProxy(target),
		cookieName: cookieName,
		realm:      realm,
		tokens:     tokens,
	}
}

// ServeHTTP implements net/http Handler
func (r *Reverse) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(r.cookieName)
	println("/cookieName", r.cookieName, "/cookie", cookie, "/err", err)
	if cookie == nil || err != nil {
		println("no cookie")
		w.Header().Add("WWW-Authenticate", r.realm)
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "Not Authorized")
		return
	}

	valid, err := r.tokens.Contains(cookie.Value)
	if err != nil {
		// TODO: think about how to log this (or if should log at all)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Oopsie!")
		return
	}
	if !valid {
		println("invalid cookie, value is", cookie.Value)
		w.Header().Add("WWW-Authenticate", r.realm)
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "Not Authorized")
		return
	}
	// finally, let the request go to the next handler
	r.actual.ServeHTTP(w, req)
}
