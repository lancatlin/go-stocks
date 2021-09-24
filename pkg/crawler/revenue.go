package crawler

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lancatlin/go-stocks/pkg/model"
)

func (c Crawler) AddRevenue(id string) {
	if err := c.updateRevenue(id); err != nil {
		panic(err)
	}
}

func (c Crawler) updateRevenue(id string) (err error) {
	fmt.Printf("%s revenue expire, crawling...\n", id)
	revenue, err := c.crawlRevenue(id)
	if err != nil {
		return
	}
	fmt.Printf("%s revenue crawled, %v\n", id, revenue)
	if err := c.Save(&revenue).Error; err != nil {
		panic(err)
	}
	c.updateRecord(model.TypeRevenue, id)
	return nil
}

func (c Crawler) crawlRevenue(id string) (revenue model.Revenue, err error) {
	file, err := download(fmt.Sprintf(c.Config.URL.Revenue, id))
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(file))
	if err != nil {
		panic(err)
	}
	s := doc.Find(`li[class="List(n)"]`).First()
	year := s.Find(`div[class="W(65px) Ta(start)"]`).Text()
	if year == "" {
		return
	}
	fmt.Println(year, ": ")
	revenue = parseRevenue(s, id)
	revenue.Time, err = time.Parse("2006/01", year)
	return revenue, err
}

func parseRevenue(s *goquery.Selection, id string) (model.Revenue) {
	revenue := model.Revenue{
		StockID: id,
	}

	s.Find("span").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 1:
			revenue.MonthRevenue = parseFloat(s.Text())
		case 3:
			revenue.YearRevenue = parseFloat(s.Text())
		}
	})
	return revenue
}
