package crawler

import (
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"strconv"
	"strings"
	"github.com/lancatlin/go-stocks/pkg/model"
)

type Crawler struct {
	db *gorm.DB
}

func New(db *gorm.DB) Crawler {
	return Crawler{db}
}


func (c Crawler) UpdatePrices() {
	file := download("http://www.twse.com.tw/exchangeReport/STOCK_DAY_ALL?response=open_data")
	if err := c.importToDatabase(file); err != nil {
		panic(err)
	}
}

func download(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(resp)
	}
	return resp.Body
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
	if c.db.NewRecord(obj) {
		method = c.db.Create
	} else {
		method = c.db.Save
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

func openDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "db.sqlite")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.Stock{})
	return db
}

