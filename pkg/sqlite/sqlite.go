package sqlite

import (
	"log"

	"github.com/iandjx/go-oauth-2/pkg/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func New() (*Client, error) {
	db, err := gorm.Open(sqlite.Open("test.db"))

	if err != nil {
		log.Fatal("could not connect database")
	}

	if err = db.AutoMigrate(&core.Client{}, &core.User{}, &core.RedirectURL{}, &core.Token{}, &core.Scope{}); err != nil {
		return nil, err
	}
	// Migration to create tables for Order and Item schema
	// db.AutoMigrate(&model.Gang, &model.Gangster{}, &model.Business{})
	return &Client{db}, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
