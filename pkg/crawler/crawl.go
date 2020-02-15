package crawler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lancatlin/go-stocks/pkg/model"
)

var (
	ErrStockNotFound = errors.New("stock not found")
)

type Crawler struct {
	*gorm.DB
}

func New(db *gorm.DB) Crawler {
	return Crawler{db}
}

func (c Crawler) isExpire() bool {
	var last model.Record
	err := c.Where("type = ? && expire_at > CURTIME()", model.TypePrice).First(&last).Error
	if gorm.IsRecordNotFoundError(err) {
		return true
	} else if err != nil {
		panic(err)
	}
	fmt.Println(last)
	return false
}

func (c Crawler) UpdateInfo() (err error) {
	if !c.isExpire() {
		return
	}
	fmt.Println("Start crawling")
	if err = c.importToDatabase("http://www.twse.com.tw/exchangeReport/STOCK_DAY_ALL?response=open_data", parseStockListed); err != nil {
		return
	}
	if err = c.importToDatabase("http://www.tpex.org.tw/web/stock/aftertrading/DAILY_CLOSE_quotes/stk_quote_result.php?l=zh-tw&o=data", parseStockCounter); err != nil {
		return
	}
	c.updatePriceRecord()
	fmt.Println("End crawling")
	return nil
}

func download(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrStockNotFound
	}
	if resp.StatusCode != 200 {
		return nil, ErrStockNotFound
	}
	return resp.Body, err
}

func (c Crawler) importToDatabase(filename string, parse func([]string) model.Stock) (err error) {
	file, err := download(filename)
	if err != nil {
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	// ignore first line
	reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		stock := parse(record)
		c.saveStock(stock)
	}
	return nil
}

func (c Crawler) saveStock(stock model.Stock) {
	var err error
	if c.First(&model.Stock{}, "id = ?", stock.ID).RecordNotFound() {
		err = c.Create(&stock).Error
	} else {
		err = c.Model(&stock).Updates(stock).Error
	}
	if err != nil {
		fmt.Println(stock)
		panic(err)
	}
}

func parseStockListed(record []string) (stock model.Stock) {
	stock.ID = record[0]
	stock.Name = record[1]
	var err error
	price := strings.Replace(record[7], ",", "", -1)
	stock.Price, err = strconv.ParseFloat(price, 64)
	if err != nil {
		panic(err)
	}
	return
}

func parseStockCounter(record []string) (stock model.Stock) {
	stock.ID = record[1]
	stock.Name = record[2]
	var err error
	price := strings.Replace(record[3], ",", "", -1)
	stock.Price, err = strconv.ParseFloat(price, 64)
	if err != nil {
		stock.Price = 0
	}
	return
}

func (c Crawler) parseRevenue(record []string) (stock model.Stock) {
	stock.ID = record[2]
	var err error
	stock.MonthRevenue, err = strconv.ParseFloat(record[9], 64)
	if err != nil {
		stock.MonthRevenue = 0
	}
	stock.YearRevenue, err = strconv.ParseFloat(record[12], 64)
	if err != nil {
		stock.YearRevenue = 0
	}
	return
}

func (c Crawler) updatePriceRecord() {
	now := time.Now()
	record := model.Record{
		Type:      model.TypePrice,
		UpdatedAt: now,
	}
	expire := time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, time.Local)
	if now.Hour() > 14 {
		expire = expire.AddDate(0, 0, 1)
	}
	record.ExpireAt = expire
	c.Save(&record)
}
