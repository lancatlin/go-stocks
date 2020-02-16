package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lancatlin/go-stocks/pkg/model"

	_ "github.com/joho/godotenv/autoload"
)

const (
	Debug   = "debug"
	Release = "release"
)

func (c *Config) openDB() (db *gorm.DB) {
	var err error
	switch c.Mode {
	case Release:
		err = c.setMySQL()
	default:
		err = c.setSqlite()
	}
	if err != nil {
		panic(err)
	}
	if err = c.DB.AutoMigrate(
		&model.Stock{}, &model.Dividend{}, &model.Record{}, &model.Revenue{},
	).Error; err != nil {
		panic(err)
	}
	return db
}

func (c *Config) setSqlite() (err error) {
	c.DB, err = gorm.Open("sqlite3", "/tmp/gorm.db")
	return
}

func (c *Config) setMySQL() (err error) {
	c.DB, err = gorm.Open(
		"mysql", fmt.Sprintf(
			"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
			c.Database.User, c.Database.Password, c.Database.Name,
		))
	return
}
