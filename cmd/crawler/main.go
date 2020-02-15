package main

import (
	"flag"

	"github.com/lancatlin/go-stocks/pkg/config"
	"github.com/lancatlin/go-stocks/pkg/crawler"
)

var (
	UpdatePrices   bool
	UpdateDividend bool
)

func init() {
	flag.BoolVar(&UpdatePrices, "p", false, "use -p to update prices")
	flag.BoolVar(&UpdateDividend, "d", false, "use -d to update dividend")
	flag.Parse()
}

func main() {
	conf := config.New()
	c := crawler.New(conf.DB)
	if UpdatePrices {
		c.UpdateInfo()
	}
	if UpdateDividend {
		c.UpdateDividends()
	}
}
