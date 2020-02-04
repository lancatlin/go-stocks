package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/lancatlin/go-stocks/pkg/model"
	"strconv"
)

func (c Crawler) UpdateDividend(id string) {
	for _, dividend := range crawlDividend(id) {
		fmt.Println(dividend)
		c.saveDividend(dividend)
	}
}

func (c Crawler) saveDividend(d model.Dividend) {
	var err error
	if c.First(&model.Dividend{}, "stock_id = ? && year = ?", d.StockID, d.Year).RecordNotFound() {
		err = c.Create(&d).Error
	} else {
		err = c.Model(&d).Updates(d).Error
	}
	if err != nil {
		fmt.Println(d)
		panic(err)
	}
}

func crawlDividend(id string) []model.Dividend {
	page, err := download("https://tw.stock.yahoo.com/d/s/dividend_" + id + ".html")
	if err != nil {
		return []model.Dividend{}
	}
	doc, err := goquery.NewDocumentFromReader(page)
	if err != nil {
		panic(err)
	}
	divs := make([]model.Dividend, 0, 10)
	doc.Find(`tr[bgcolor='#FFFFFF']`).Each(func(i int, s *goquery.Selection) {
		dividend := parseDividend(s, id)
		divs = append(divs, dividend)
	})
	return divs
}

func parseDividend(s *goquery.Selection, id string) model.Dividend {
	dividend := model.Dividend{
		StockID: id,
	}
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
	return dividend
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
