package sqlite

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sort"

	"github.com/iandjx/go-oauth-2/pkg/core"
	"github.com/iandjx/go-oauth-2/pkg/util"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db     *gorm.DB
	Secret string
}

func NewClientRepository(c *Client, secret string) *ClientRepository {
	return &ClientRepository{c.db, secret}
}

func (c *ClientRepository) CreateClient(name string, redirectURLS []string, scope []string) (*core.Client, error) {
	token := randToken()
	var formattedURLs []core.RedirectURL
	for _, value := range redirectURLS {
		url := core.RedirectURL{URL: value}
		formattedURLs = append(formattedURLs, url)
	}
	var scopes []core.Scope
	for _, value := range scope {
		access := core.Scope{Access: value}
		scopes = append(scopes, access)
	}

	client := core.Client{Name: name, RedirectURLs: formattedURLs, Scopes: scopes, ClientSecret: token}
	tx := c.db.Create(&client)
	if tx.Error != nil {
		return nil, tx.Error
	}
	et, err := util.Encrypt(token, c.Secret)
	if err != nil {
		return nil, err
	}
	client.ClientSecret = et

	return &client, nil
}

func (c *ClientRepository) GetClient(clientID uint, redirectURL string, scope []string) (*core.Client, error) {
	var client core.Client

	tx := c.db.First(&client, clientID)

	if err := tx.Error; err != nil {
		return nil, err
	}

	_, redirectURLFound := Find(client.RedirectURLs, redirectURL)

	if !redirectURLFound {
		return nil, errors.New("invalid redirect url")
	}
	var scopes []core.Scope

	for _, v := range scope {
		ns := core.Scope{Access: v}
		scopes = append(scopes, ns)
	}

	sort.Slice(scopes, func(i, j int) bool {
		return scopes[i].Access < scopes[j].Access
	})

	sort.Slice(client.Scopes, func(i, j int) bool {
		return scopes[i].Access < scopes[j].Access
	})
	if !equalScope(scopes, client.Scopes) {
		return nil, errors.New("invalid scope")
	}

	client.ClientSecret = ""
	return &client, nil

}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []core.RedirectURL, val string) (int, bool) {
	for i, item := range slice {
		if item.URL == val {
			return i, true
		}
	}
	return -1, false
}

func randToken() string {
	b := make([]byte, 20)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func equalScope(a, b []core.Scope) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.Access != b[i].Access {
			return false
		}
	}
	return true
}
