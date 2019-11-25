package main

import (
	"net/http"
	"io"
	"encoding/csv"
	"strconv"
	"github.com/jinzhu/gorm"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"strings"
)

func downloadCSV(url string) io.Reader {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	return resp.Body
}

type Stock struct {
	Number string `gorm:"primary_key"`
	Name string
	Price float64
}

func importToDatabase(db *gorm.DB, file io.Reader) (err error) {
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
	stock.Number = record[0]
	stock.Name =  record[1]
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
	file := downloadCSV("http://www.twse.com.tw/exchangeReport/STOCK_DAY_ALL?response=open_data")
	if err := importToDatabase(db, file); err != nil {
		panic(err)
	}
}
