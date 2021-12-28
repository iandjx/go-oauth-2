package sqlite

import (
	"github.com/iandjx/go-oauth-2/pkg/core"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(c *Client) *UserRepository {
	return &UserRepository{c.db}
}

func (u *UserRepository) CreateUser(name string, email string, password string) (*core.User, error) {
	user := core.User{Name: name, Email: email, Password: password}
	tx := u.db.Create(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (u *UserRepository) GetUserByID(userID uint) (*core.User, error) {
	var user core.User

	tx := u.db.First(&user, userID)

	if err := tx.Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetUserByCredentials(email string, password string) (*core.User, error) {
	var user core.User

	tx := u.db.First(&user, "name = ? AND password = ?", email, password)

	if err := tx.Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) AddToken(userID uint, tokenType string, token string) error {
	var user core.User
	tx := u.db.First(&user, userID)

	if err := tx.Error; err != nil {
		return err
	}
	user.Tokens = append(user.Tokens, core.Token{Type: tokenType, Token: token})

	tx2 := u.db.Save(&user)

	if err := tx2.Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) DeleteToken(userID uint, token string) error {
	var user core.User
	tx := u.db.First(&user, userID)
	if err := tx.Error; err != nil {
		return err
	}
	user.Tokens = findAndDeleteToken(user.Tokens, token)

	tx2 := u.db.Save(&user)

	if err := tx2.Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) DeleteAllTokens(userID uint) error {
	var user core.User
	tx := u.db.First(&user, userID)
	if err := tx.Error; err != nil {
		return err
	}
	user.Tokens = []core.Token{}

	tx2 := u.db.Save(&user)

	if err := tx2.Error; err != nil {
		return err
	}

	return nil
}
func findAndDeleteToken(tokens []core.Token, token string) []core.Token {
	index := 0
	for _, i := range tokens {
		if i.Token != token {
			tokens[index] = i
			index++
		}
	}
	return tokens[:index]
}
