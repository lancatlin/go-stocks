package crawler

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lancatlin/go-stocks/pkg/model"
)

func (c Crawler) UpdateRevenues() (err error) {
	stocks := c.findStocks("revenues")
	fmt.Println(stocks)
	for _, stockID := range stocks {
		c.updateRevenue(stockID)
	}
	return nil
}

func (c Crawler) AddRevenue(id string) {
	if err := c.updateRevenue(id); err != nil {
		panic(err)
	}
}

func (c Crawler) updateRevenue(id string) (err error) {
	if !c.isExpire(model.TypeRevenue, id) {
		return nil
	}
	revenue := c.crawlRevenue(id)
	if same, hash := c.isSame(revenue, model.TypeRevenue, id); !same {
		c.Save(&revenue)
		c.updateRecord(model.TypeRevenue, id, hash)
	}
	return nil
}

func (c Crawler) crawlRevenue(id string) (revenue model.Revenue) {
	file, err := download(fmt.Sprintf(c.Config.URL.Revenue, id))
	if err != nil {
		return model.Revenue{}
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(file))
	if err != nil {
		panic(err)
	}
	doc.Find(`tr[bgcolor='#FFFFFF']`).Slice(2, 14).EachWithBreak(func(i int, s *goquery.Selection) bool {
		if r, ok := parseRevenue(i, s, id); ok {
			revenue = r
		} else {
			return false
		}
		return true
	})
	return revenue
}

func parseRevenue(month int, s *goquery.Selection, id string) (model.Revenue, bool) {
	revenue := model.Revenue{
		StockID: id,
	}

	now := time.Now()
	if now.Month() == time.January {
		revenue.Time = time.Date(now.Year()-1, time.Month(month+1), 1, 0, 0, 0, 0, time.Local)
	} else {
		revenue.Time = time.Date(now.Year(), time.Month(month+1), 1, 0, 0, 0, 0, time.Local)
	}

	c := s.Children()
	m := c.Get(5).FirstChild.Data
	y := c.Get(7).FirstChild.Data
	if m == "-" || y == "-" {
		return revenue, false
	}

	revenue.MonthRevenue = parseFloat(m)
	revenue.YearRevenue = parseFloat(y)
	fmt.Println(revenue)
	return revenue, true
}
