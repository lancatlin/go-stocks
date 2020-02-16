package crawler

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
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
	if !c.isRevExpire(id) {
		return nil
	}
	revenue := c.crawlRevenue(id)
	if same, hash := c.isRevSame(revenue); !same {
		c.Save(&revenue)
		c.updateRevRecord(id, hash)
	}
	return nil
}

func (c Crawler) isRevExpire(id string) bool {
	err := c.Where("type = ? and stock_id = ? and expire_at > ?", model.TypeRevenue, id, time.Now()).First(&model.Record{}).Error
	if gorm.IsRecordNotFoundError(err) {
		return true
	} else if err != nil {
		panic(err)
	}
	return false
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

func (c Crawler) isRevSame(revenue model.Revenue) (bool, string) {
	hash := hashString(revenue)
	var last model.Record
	err := c.Where("type = ? and stock_id = ?", model.TypeRevenue, revenue.StockID).Order("updated_at desc").First(&last).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, hash
	} else if err != nil {
		panic(err)
	}
	return hash == last.Hash, hash
}

func (c Crawler) updateRevRecord(id, hash string) {
	now := time.Now()
	record := model.Record{
		Type:      model.TypeRevenue,
		StockID:   id,
		Hash:      hash,
		UpdatedAt: now,
	}
	record.ExpireAt = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local)
	if err := c.Create(&record).Error; err != nil {
		panic(err)
	}
}
