package client

import (
	"context"

	"github.com/iandjx/go-oauth-2/pkg/auth"
	"github.com/iandjx/go-oauth-2/pkg/core"
)

type Service interface {
	CreateClient(ctx context.Context, p CreateParam) (*core.Client, error)
}

type service struct {
	clientRepo core.ClientRepository
	JWTSecret  string
}

func NewService(cr core.ClientRepository, JWTSecret string) Service {
	return &service{cr, JWTSecret}
}

func (s *service) CreateClient(ctx context.Context, p CreateParam) (*core.Client, error) {
	au := auth.FromContext(ctx)
	if au == nil || au.ID == 0 {
		err := auth.ErrAuthRequired
		return nil, err
	}
	c, err := s.clientRepo.CreateClient(p.Name, p.RedirectURLs, p.Scope)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) GetClient(ctx context.Context, clientID uint, redirectURL string, scope []string) (*core.Client, error) {

	au := auth.FromContext(ctx)
	if au == nil || au.ID == 0 {
		err := auth.ErrAuthRequired
		return nil, err
	}
	c, err := s.clientRepo.GetClient(clientID, redirectURL, scope)
	if err != nil {
		return nil, err
	}
	c.ClientSecret = ""
	return c, nil
}

// TODO create func to delete client
