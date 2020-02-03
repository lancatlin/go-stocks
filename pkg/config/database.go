package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lancatlin/go-stocks/pkg/model"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Mode string

const (
	DebugMode   Mode = "debug mode"
	ReleaseMode Mode = "release mode"
)

func getUpdateTime() time.Duration {
	n, err := strconv.Atoi(os.Getenv("UPDATE"))
	if err != nil {
		n = 10
	}
	return time.Minute * time.Duration(n)
}

type Config struct {
	Mode   Mode
	DB     *gorm.DB
	Update time.Duration
	Host   string
	Port   string
	Base   string
}

func New() Config {
	mode := ReleaseMode
	config := Config{
		Mode:   mode,
		DB:     openDB(mode),
		Update: getUpdateTime(),
		Host:   os.Getenv("HOST"),
		Port:   os.Getenv("PORT"),
		Base:   os.Getenv("BASE"),
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
