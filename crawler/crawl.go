package main

import (
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io"
	"net/http"
	"strconv"
	"strings"
)

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

type Stock struct {
	ID        string `gorm:"primary_key"`
	Name      string
	Price     float64
	Dividends []Dividend
}

func importToDatabase(db *gorm.DB, file io.ReadCloser) (err error) {
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
		var method func(interface{}) *gorm.DB
		if db.NewRecord(&stock) {
			method = db.Create
		} else {
			method = db.Save
		}
		if err := method(&stock).Error; err != nil {
			return err
		}
	}
	return nil
}

func parseStock(record []string) (stock Stock) {
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
	db.AutoMigrate(&Stock{})
	return db
}

func main() {
	db := openDB()
	file := download("http://www.twse.com.tw/exchangeReport/STOCK_DAY_ALL?response=open_data")
	if err := importToDatabase(db, file); err != nil {
		panic(err)
	}
}
