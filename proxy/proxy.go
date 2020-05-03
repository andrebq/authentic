package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/andrebq/authentic/internal/session"
	"github.com/andrebq/authentic/internal/webflow"
)

type (
	// Reverse proxy with cookie checks
	Reverse struct {
		actual     *httputil.ReverseProxy
		cookieName string
		realm      string
		s          *session.S
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
func NewReverse(cookieName, realm string, session *session.S, target *url.URL) *Reverse {
	return &Reverse{
		actual:     httputil.NewSingleHostReverseProxy(target),
		cookieName: cookieName,
		realm:      realm,
		s:          session,
	}
}

// ServeHTTP implements net/http Handler
func (r *Reverse) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(r.cookieName)
	if cookie == nil || err != nil {
		webflow.Authenticate(w, req, r.realm, "Please authenticate")
		return
	}

	err = r.s.Verify(cookie.Value)
	switch {
	case session.IsExpired(err):
		webflow.Authenticate(w, req, r.realm, "Login expired, please re-authenticate.")
		return
	case session.IsInvalid(err):
		webflow.Authenticate(w, req, r.realm, "Invalid credentials, please authenticate.")
		return
	case err != nil:
		webflow.InternalError(w, req)
		return
	}
	r.actual.ServeHTTP(w, req)
}
