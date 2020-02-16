package crawler

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lancatlin/go-stocks/pkg/model"
)

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
	fmt.Printf("Month is %d\n", month)
	fmt.Println(s.Text())
	now := time.Now()
	if now.Month() == time.January {
		revenue.Time = time.Date(now.Year()-1, time.Month(month+1), 1, 0, 0, 0, 0, time.Local)
	} else {
		revenue.Time = time.Date(now.Year(), time.Month(month+1), 1, 0, 0, 0, 0, time.Local)
	}
	ok := true
	s.Children().EachWithBreak(func(i int, s *goquery.Selection) bool {
		switch i {
		case 5:
			// If data is unvalid, return false
			if s.Text() == "-" {
				ok = false
				return false
			}
			revenue.MonthRevenue = parseFloat(s.Text())
		case 7:
			if s.Text() == "-" {
				ok = false
				return false
			}
			revenue.YearRevenue = parseFloat(s.Text())
		}
		return true
	})
	fmt.Println(revenue)
	return revenue, ok
}
