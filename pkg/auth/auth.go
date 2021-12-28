package auth

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrAuthRequired = errors.New("authorization required")
	ErrAuthAccess   = errors.New("access denied")
)

// User represents authorized user
type User struct {
	ID    uint
	Email string
}

type ctxKey int

const authKey ctxKey = iota

// ToContext set user details to context.
func ToContext(parent context.Context, u *User) context.Context {
	return context.WithValue(parent, authKey, u)
}

// FromContext returns a user details from the given context if one is present.
//
// Returns nil if usser detail cannot be found.
func FromContext(ctx context.Context) *User {
	if ctx == nil {
		return nil
	}
	if u, ok := ctx.Value(authKey).(*User); ok {
		return u
	}
	return nil
}

// FromRequest returns a user details from the given request if one is present.
//
// Returns nil if usser detail cannot be found.
func FromRequest(r *http.Request) *User {
	ctx := r.Context()
	return FromContext(ctx)
}
