package firebase

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/dghubble/sling"
)

type (
	// UserCatalog exposes some parts of Firebase user auth
	UserCatalog struct {
		fbApiKey string
	}
)

// Users creates a new user catalog
func Users() (*UserCatalog, error) {
	uc := &UserCatalog{
		fbApiKey: os.Getenv("FIREBASE_WEB_APIKEY"),
	}
	if uc.fbApiKey == "" {
		return nil, errors.New("Missing FIREBASE_WEB_APIKEY environment variable")
	}
	return uc, nil
}

// Authenticate checks if the combination username/password are valid
func (uc *UserCatalog) Authenticate(ctx context.Context, username, password string) error {
	qs := struct {
		APIKey string `url:"key"`
	}{APIKey: uc.fbApiKey}
	signIn := struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		SecureToken bool   `json:"returnSecureToken"`
	}{
		Email:       username,
		Password:    password,
		SecureToken: true,
	}
	println("info: ", fmt.Sprintf("%v", signIn))
	response := struct {
		IDToken    string `json:"idToken"`
		Email      string `json:"email"`
		Registered bool   `json:"registered"`
	}{}
	httpResponse, err := sling.New().Post("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword").
		QueryStruct(qs).BodyJSON(signIn).Receive(&response, nil)
	if err != nil {
		println("network error")
		return err
	}
	if httpResponse.StatusCode != 200 {
		println("status code error", httpResponse.StatusCode)
		return errors.New("unable to login")
	}
	return nil
}
