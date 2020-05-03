package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type (
	// UserCatalog exposes some parts of Firebase user auth
	UserCatalog struct {
		app  *firebase.App
		auth *auth.Client
	}
)

// Users creates a new user catalog
func Users(ctx context.Context) (*UserCatalog, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to open user catalog: %v", err)
	}
	auth, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to open user catalog: %v", err)
	}
	return &UserCatalog{app, auth}, nil
}

// Authenticate checks if the combination username/password are valid
func (uc *UserCatalog) Authenticate(ctx context.Context, username, password string) error {
	tk, err := uc.auth.VerifyIDTokenAndCheckRevoked(ctx, password)
	if err != nil {
		return fmt.Errorf("unable to validate credentials: %v", err)
	}
	user, err := uc.auth.GetUser(ctx, tk.UID)
	if err != nil {
		return fmt.Errorf("unable to validate credentials: %v", err)
	}
	if user.Email != username {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}
