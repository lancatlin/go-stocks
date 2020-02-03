package crawler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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

func (c Crawler) UpdateInfo() (err error) {
	fmt.Println("Start crawling")
	if err = c.importToDatabase("http://www.twse.com.tw/exchangeReport/STOCK_DAY_ALL?response=open_data", parseStockListed); err != nil {
		return
	}
	if err = c.importToDatabase("http://www.tpex.org.tw/web/stock/aftertrading/DAILY_CLOSE_quotes/stk_quote_result.php?l=zh-tw&o=data", parseStockCounter); err != nil {
		return
	}
	if err = c.importToDatabase("http://mopsfin.twse.com.tw/opendata/t187ap05_L.csv", c.parseRevenue); err != nil {
		return
	}
	if err = c.importToDatabase("http://mopsfin.twse.com.tw/opendata/t187ap05_O.csv", c.parseRevenue); err != nil {
		return
	}
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
		c.save(&stock)
	}
	return nil
}

func (c Crawler) save(obj interface{}) {
	var method func(interface{}) *gorm.DB
	if c.NewRecord(obj) {
		method = c.Create
	} else {
		method = func(obj interface{}) *gorm.DB { return c.First(obj).Updates(obj) }
	}
	if err := method(obj).Error; err != nil {
		fmt.Println(obj)
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
	stock.CompareLastYear, err = strconv.ParseFloat(record[9], 64)
	if err != nil {
		stock.CompareLastYear = 0
	}
	return
}
