package sqlite

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/iandjx/go-oauth-2/pkg/core"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db     *gorm.DB
	Secret string
}

func NewClientRepository(c *Client, secret string) *ClientRepository {
	return &ClientRepository{c.db, secret}
}

func (c *ClientRepository) CreateClient(name string, redirectURLSs []string, scope []string) (*core.Client, error) {
	token := randToken()
	client := core.Client{Name: name, RedirectURLs: redirectURLSs, Scope: scope, ClientSecret: token}
	tx := c.db.Create(&client)
	if tx.Error != nil {
		return nil, tx.Error
	}
	et, err := encrypt(token, c.Secret)
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

	if !cmp.Equal(scope, client.Scope) {
		return nil, errors.New("invalid scope")
	}

	return &client, nil

}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
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

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Encrypt method is to encrypt or hide any classified text
func encrypt(text, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

func decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Decrypt method is to extract back the encrypted text
func decrypt(text, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
