package main

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

type Dividend struct {
	Stock         Stock  `gorm:"foreignkey:StockID"`
	StockID       string `gorm:"primary_key"`
	Year          int    `gorm:"primary_key"`
	MoneyDividend float64
	StockDividend float64
}

func CrawlDividend(id string) []Dividend {
	page := download("https://tw.stock.yahoo.com/d/s/dividend_" + id + ".html")
	doc, err := goquery.NewDocumentFromReader(page)
	if err != nil {
		panic(err)
	}
	divs := make([]Dividend, 10)
	doc.Find(`tr[bgcolor='#FFFFFF']`).Each(func(i int, s *goquery.Selection) {
		dividend := Dividend{
			StockID: id,
		}
		dividend.parseDividend(s)
		divs[i] = dividend
	})
	return divs
}

func (dividend *Dividend) parseDividend(s *goquery.Selection) {
	s.Children().Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			dividend.Year = parseInt(s.Text())
		case 1:
			dividend.MoneyDividend = parseFloat(s.Text())
		case 4:
			dividend.StockDividend = parseFloat(s.Text())
		}
	})
}

func parseFloat(s string) float64 {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return num
}

func parseInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return num
}
