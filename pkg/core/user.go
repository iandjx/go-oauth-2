package core

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Tokens   []Token
}

type Token struct {
	gorm.Model
	UserID uint
	User   User
	Type   string
	Token  string
}

type UserRepository interface {
	CreateUser(name string, email string, password string) (*User, error)
	GetUserByID(userID uint) (*User, error)
	GetUserByCredentials(email string, password string) (*User, error)
	AddToken(userID uint, tokenType string, token string) error
	DeleteToken(userID uint, token string) error
	DeleteAllTokens(userID uint) error
}
