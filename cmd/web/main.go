package main

import (
	"github.com/lancatlin/go-stocks/pkg/config"
	"github.com/lancatlin/go-stocks/pkg/web"
)

func main() {
	conf := config.New()
	router := web.Registry(conf)
	router.Run(":8080")
}
