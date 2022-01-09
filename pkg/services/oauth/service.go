package oauth

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/iandjx/go-oauth-2/pkg/auth"
	"github.com/iandjx/go-oauth-2/pkg/core"
)

type Service interface {
	GenerateToken(ctx context.Context, authCode string, clientID uint, clientSecret string, scope []string, redirectURL string) (string, error)
	GenerateAuthCode(ctx context.Context, clientID uint, redirectURL string, clientSecret string, scope []string) (string, error)
}

type service struct {
	userRepo   core.UserRepository
	clientRepo core.ClientRepository
	JWTSecret  string
}

type OAuthClaim struct {
	userID       uint
	clientID     uint
	clientSecret string
	scope        []string
	jwt.StandardClaims
}

func NewService(ur core.UserRepository, cr core.ClientRepository, JWTSecret string) Service {
	return &service{ur, cr, JWTSecret}
}

func (s *service) verifyClient(clientID uint, redirectURL string, scope []string) error {
	_, err := s.clientRepo.GetClient(clientID, redirectURL, scope)

	if err != nil {
		return err
	}
	return nil
}

func (s *service) GenerateToken(ctx context.Context, authCode string, clientID uint, clientSecret string, scope []string, redirectURL string) (string, error) {
	au := auth.FromContext(ctx)
	if au == nil || au.ID == 0 {
		err := auth.ErrAuthRequired
		return "", err
	}

	err := s.verifyClient(clientID, redirectURL, scope)
	if err != nil {
		return "", err
	}
	err = s.verifyAuthCode(authCode, clientID, clientSecret, scope)

	if err != nil {
		return "", err
	}
	s.userRepo.DeleteToken(uint(au.ID), authCode)

	tokenType := "access_token"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, OAuthClaim{
		clientID:     clientID,
		userID:       au.ID,
		clientSecret: clientSecret,
		scope:        scope,
	})

	tokenString, err := token.SignedString([]byte(s.JWTSecret))

	if err != nil {
		return "", err
	}

	err = s.userRepo.AddToken(au.ID, tokenType, tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *service) GenerateAuthCode(ctx context.Context, clientID uint, redirectURL string, clientSecret string, scope []string) (string, error) {
	au := auth.FromContext(ctx)
	if au == nil || au.ID == 0 {
		err := auth.ErrAuthRequired
		return "", err
	}
	err := s.verifyClient(clientID, redirectURL, scope)

	if err != nil {
		return "", err
	}

	tokenType := "oauth"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, OAuthClaim{
		clientID:     clientID,
		userID:       au.ID,
		clientSecret: clientSecret,
		scope:        scope,
	})

	tokenString, err := token.SignedString([]byte(s.JWTSecret))

	if err != nil {
		return "", err
	}

	err = s.userRepo.AddToken(au.ID, tokenType, tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *service) verifyAuthCode(authCode string, clientId uint, clientSecret string, scope []string) error {

	var decodedClientID uint
	var decodedUserID uint
	var decodedClientSecret string
	var decodedScope []string

	token, err := jwt.ParseWithClaims(authCode, OAuthClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.JWTSecret), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(OAuthClaim); ok && token.Valid {
		decodedClientID = claims.clientID
		decodedUserID = claims.userID
		decodedClientSecret = claims.clientSecret
		decodedScope = claims.scope
	} else {
		return errors.New("invalid token")
	}

	if decodedClientID != clientId || decodedClientSecret != clientSecret || !cmp.Equal(scope, decodedScope) {
		return errors.New("code does not belong to project")
	}

	_, err = s.userRepo.GetUserByID(decodedUserID)
	if err != nil {
		return err
	}
	return nil
}
