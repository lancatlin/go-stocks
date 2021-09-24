package crawler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lancatlin/go-stocks/pkg/model"
)

func (c Crawler) UpdateDividend(id string) (err error) {
	divs := c.crawlDividend(id)
	for _, dividend := range divs {
		fmt.Println(dividend)
		if err := c.saveDividend(dividend); err != nil {
			panic(err)
		}
	}
	c.updateRecord(model.TypeDividend, id)
	return nil
}

func (c Crawler) saveDividend(d model.Dividend) (err error) {
	if c.First(&model.Dividend{}, "stock_id = ? and year = ?", d.StockID, d.Year).RecordNotFound() {
		err = c.Create(&d).Error
	} else {
		err = c.Model(&d).Updates(d).Error
	}
	return
}

func (c Crawler) crawlDividend(id string) []model.Dividend {
	file, err := download(fmt.Sprintf(c.Config.URL.Dividend, id))
	if err != nil {
		return []model.Dividend{}
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(file))
	if err != nil {
		panic(err)
	}
	divs := make([]model.Dividend, 0, 10)
	doc.Find(`li[class="List(n)"]`).Each(func(i int, s *goquery.Selection) {
		year := s.Find(`div[class="D(f) W(98px) Ta(start)"]`).Text()
		if year == "" {
			return
		}
		dividend := parseDividend(s, id)
		dividend.Year = parseInt(year)
		divs = append(divs, dividend)
	})
	return divs
}

func parseDividend(s *goquery.Selection, id string) model.Dividend {
	dividend := model.Dividend{
		StockID: id,
	}
	s.Find("span").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			dividend.MoneyDividend = parseFloat(s.Text())
		case 1:
			dividend.StockDividend = parseFloat(s.Text())
		}
	})
	return dividend
}

func parseFloat(s string) float64 {
	re := regexp.MustCompile(`-?[\d\.]+`)
	num, err := strconv.ParseFloat(re.FindString(s), 64)
	if err != nil {
		return 0
	}
	return num
}

func parseInt(s string) int {
	re := regexp.MustCompile(`^\d+`)
	num, err := strconv.Atoi(re.FindString(s))
	if err != nil {
		panic(err)
	}
	return num
}
