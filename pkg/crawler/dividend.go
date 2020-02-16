package crawler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/lancatlin/go-stocks/pkg/model"
)

func (c Crawler) UpdateDividends() (err error) {
	stocks := c.findStocks("dividends")
	fmt.Println(stocks)
	for _, stockID := range stocks {
		if c.isDividendExpire(stockID) {
			c.UpdateDividend(stockID)
		}
	}
	return nil
}

func (c Crawler) findStocks(table string) (stocks []string) {
	rows, err := c.Table(table).Select("stock_id").Group("stock_id").Rows()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var stock string
		if err := rows.Scan(&stock); err != nil {
			panic(err)
		}
		stocks = append(stocks, stock)
	}
	return
}

func (c Crawler) isDividendExpire(stock string) bool {
	var last model.Record
	err := c.Where("type = ? and stock_id = ? and expire_at > ?", model.TypeDividend, stock, time.Now()).First(&last).Error
	if gorm.IsRecordNotFoundError(err) {
		return true
	} else if err != nil {
		panic(err)
	}
	fmt.Println(last)
	return false
}

func (c Crawler) UpdateDividend(id string) {
	divs := c.crawlDividend(id)
	if same, hash := c.isDivSame(id, divs); !same {
		for _, dividend := range divs {
			fmt.Println(dividend)
			if err := c.saveDividend(dividend); err != nil {
				panic(err)
			}
		}
		c.updateDividendRecord(id, hash)
	} else {
		fmt.Printf("%s not change %s\n", id, hash)
	}
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
	re := regexp.MustCompile(`-?[\d\.]+`)
	num, err := strconv.ParseFloat(re.FindString(s), 64)
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

func (c Crawler) updateDividendRecord(stockID, hash string) {
	now := time.Now()
	record := model.Record{
		Type:      model.TypeDividend,
		StockID:   stockID,
		Hash:      hash,
		UpdatedAt: now,
	}
	expire := time.Date(now.Year(), time.June, 1, 0, 0, 0, 0, time.Local)
	if now.After(expire) {
		expire = expire.AddDate(1, 0, 0)
	}
	record.ExpireAt = expire
	c.Save(&record)
}
