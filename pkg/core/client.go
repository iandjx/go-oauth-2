package core

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ClientSecret string
	Name         string
	RedirectURLs []RedirectURL
	Scopes       []Scope
}

type RedirectURL struct {
	gorm.Model
	Client   Client
	ClientID uint
	URL      string
}

type Scope struct {
	gorm.Model
	Client   Client
	ClientID uint
	Access   string
}

type ClientRepository interface {
	CreateClient(name string, redirecURLs []string, scope []string) (*Client, error)
	GetClient(clientID uint, redirectURL string, scope []string) (*Client, error)
}
