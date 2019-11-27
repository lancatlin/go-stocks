package main

import (
	"flag"
	"github.com/lancatlin/go-stocks/pkg/config"
	"github.com/lancatlin/go-stocks/pkg/crawler"
)

var (
	UpdatePrices   bool
	UpdateDividend string
)

func init() {
	flag.BoolVar(&UpdatePrices, "p", false, "use -p to update prices")
	flag.StringVar(&UpdateDividend, "d", "", "use -d 'stock_number' to update dividend")
	flag.Parse()
}

func main() {
	conf := config.New()
	c := crawler.New(conf.DB)
	if UpdatePrices {
		c.UpdatePrices()
	}
	if UpdateDividend != "" {
		c.UpdateDividend(UpdateDividend)
	}
}
