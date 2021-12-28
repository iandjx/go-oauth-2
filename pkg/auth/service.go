package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	UserFrom(r *http.Request) (*User, error)
}

type AuthClaim struct {
	ID    uint
	Email string
	jwt.StandardClaims
}

type service struct {
	JWTSecret string
}

func NewService(JWTSecret string) AuthService {
	return &service{JWTSecret}
}

const (
	authToken = "x-auth"
)

func (s *service) UserFrom(r *http.Request) (*User, error) {
	hc, _ := r.Cookie(authToken)
	if hc == nil {
		return nil, errors.New("missing token from request")
	}
	return s.getUserFromToken(hc.Value)
}

func (s *service) getUserFromToken(token string) (*User, error) {
	t, err := jwt.ParseWithClaims(token, &AuthClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	ac, ok := t.Claims.(*AuthClaim)
	if !ok {
		return nil, fmt.Errorf("could not parse claims")
	}

	u := &User{
		ID:    ac.ID,
		Email: ac.Email,
	}
	return u, nil
}
