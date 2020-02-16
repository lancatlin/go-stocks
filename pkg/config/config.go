package config

import (
	"os"

	"github.com/jinzhu/gorm"
	"github.com/naoina/toml"
)

type Config struct {
	Mode   string
	DB     *gorm.DB
	Update int64

	Server struct {
		Host string
		Port int
		Base string
	}

	Database struct {
		Name     string
		User     string
		Password string
	}

	URL struct {
		Listed   string
		Counter  string
		Dividend string
		Revenue  string
	}
}

func New() Config {
	f, err := os.Open("config.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var config Config
	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}
	config.openDB()
	return config
}
