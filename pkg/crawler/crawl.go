package crawler

import (
	"encoding/csv"
	"errors"
	"fmt"
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
	if err = c.importToDatabase(file); err != nil {
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

func (c Crawler) importToDatabase(file io.ReadCloser) (err error) {
	defer file.Close()
	reader := csv.NewReader(file)
	// ignore first line
	reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		stock := parseStock(record)
		fmt.Println(stock)
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

func parseStock(record []string) (stock model.Stock) {
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
