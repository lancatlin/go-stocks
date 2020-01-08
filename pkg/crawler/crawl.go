package crawler

import (
	"encoding/csv"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/lancatlin/go-stocks/pkg/model"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func (c Crawler) UpdatePrices() (err error) {
	file, err := download("http://www.twse.com.tw/exchangeReport/STOCK_DAY_ALL?response=open_data")
	if err != nil {
		return
	}
	if err = c.importToDatabase(file, parseStockListed); err != nil {
		return
	}
	file, err = download("http://www.tpex.org.tw/web/stock/aftertrading/DAILY_CLOSE_quotes/stk_quote_result.php?l=zh-tw&o=data")
	if err != nil {
		return
	}
	if err = c.importToDatabase(file, parseStockCounter); err != nil {
		return
	}
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

func (c Crawler) importToDatabase(file io.ReadCloser, parse func([]string) model.Stock) (err error) {
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
		// fmt.Println(stock)
		c.save(&stock)
	}
	return nil
}

func (c Crawler) save(obj interface{}) {
	var method func(interface{}) *gorm.DB
	if c.NewRecord(obj) {
		method = c.Create
	} else {
		method = c.Save
	}
	if err := method(obj).Error; err != nil {
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
