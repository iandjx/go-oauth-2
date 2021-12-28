package user

import (
	"context"

	"github.com/golang-jwt/jwt"
	"github.com/iandjx/go-oauth-2/pkg/auth"
	"github.com/iandjx/go-oauth-2/pkg/core"
)

type Service interface {
	Register(name string, email string, password string) (string, error)
}

type service struct {
	userRepo  core.UserRepository
	JWTSecret string
}

func NewService(ur core.UserRepository, JWTSecret string) Service {
	return &service{ur, JWTSecret}
}

func (s *service) Register(name string, email string, password string) (string, error) {
	u, err := s.userRepo.CreateUser(name, email, password)
	if err != nil {
		return "", err
	}
	at, err := s.generateAuthToken(u.ID, u.Email)
	if err != nil {
		return "", err
	}
	return at, nil
}

func (s *service) Login(email string, password string) (string, error) {
	u, err := s.userRepo.GetUserByCredentials(email, password)

	if err != nil {
		return "", err
	}

	at, err := s.generateAuthToken(u.ID, u.Email)
	if err != nil {
		return "", err
	}
	return at, nil
}

func (s *service) Logout(ctx context.Context) error {
	au := auth.FromContext(ctx)
	if au == nil || au.ID == 0 {
		err := auth.ErrAuthRequired
		return err
	}

	u, err := s.userRepo.GetUserByID(au.ID)
	if err != nil {
		return err
	}
	err = s.userRepo.DeleteAllTokens(u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) generateAuthToken(userID uint, email string) (string, error) {
	tokenType := "auth"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.AuthClaim{
		ID:    userID,
		Email: email,
	})

	tokenString, err := token.SignedString([]byte(s.JWTSecret))

	if err != nil {
		return "", err
	}

	err = s.userRepo.AddToken(userID, tokenType, tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
