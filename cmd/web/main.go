package main

import (
	"fmt"

	"github.com/lancatlin/go-stocks/pkg/config"
	"github.com/lancatlin/go-stocks/pkg/web"
)

func main() {
	conf := config.New()
	router := web.Registry(conf)
	fmt.Printf("Running in %s mode", conf.Mode)
	if err := router.Run(
		fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)); err != nil {
		panic(err)
	}
}
