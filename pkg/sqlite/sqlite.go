package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func New() (*Client, error) {
	var err error
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Migration to create tables for Order and Item schema
	// db.AutoMigrate(&model.Gang, &model.Gangster{}, &model.Business{})
	return &Client{db}, nil
}

func (c *Client) MigrateSchema(dst ...interface{}) error {
	err := c.db.AutoMigrate(dst...)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
