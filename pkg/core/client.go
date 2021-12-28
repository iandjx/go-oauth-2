package core

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ClientSecret string
	Name         string
	RedirectURLs []string
	Scope        []string
}

type ClientRepository interface {
	CreateClient(name string, redirecURLs []string, scope []string) (*Client, error)
	GetClient(clientID uint, redirectURL string, scope []string) (*Client, error)
}
