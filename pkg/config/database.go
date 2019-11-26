package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lancatlin/go-stocks/pkg/model"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Mode string

const (
	DebugMode   Mode = "debug mode"
	ReleaseMode Mode = "release mode"
)

type Config struct {
	Mode Mode
	DB   *gorm.DB
	Host string
	Port string
}

func New() Config {
	mode := ReleaseMode
	config := Config{
		Mode: mode,
		DB:   openDB(mode),
	}
	return config
}

func openDB(mode Mode) (db *gorm.DB) {
	var err error
	switch mode {
	case DebugMode:
		db, err = gorm.Open("sqlite3", "/tmp/gorm.db")
	case ReleaseMode:
		conf := struct {
			Database string
			User     string
			Password string
		}{
			os.Getenv("DB"),
			os.Getenv("DB"),
			os.Getenv("PASSWORD"),
		}
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", conf.User, conf.Password, conf.Database))
	}
	if err != nil {
		panic(err)
	}
	if err = db.AutoMigrate(&model.Stock{}, &model.Dividend{}).Error; err != nil {
		panic(err)
	}
	return db
}
