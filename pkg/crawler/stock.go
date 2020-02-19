package crawler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/lancatlin/go-stocks/pkg/config"
	"github.com/lancatlin/go-stocks/pkg/model"
)

var (
	ErrStockNotFound = errors.New("stock not found")
)

type Crawler struct {
	*gorm.DB
	Config config.Config
}

func New(config config.Config) Crawler {
	return Crawler{config.DB, config}
}

func (c Crawler) UpdateInfo() (err error) {
	fmt.Println("Start crawling")
	fmt.Println("Crawl Listed")
	if err = c.updateStocks(model.TypePriceListed, parseStockListed); err != nil {
		return
	}
	fmt.Println("Crawl Counter")
	if err = c.updateStocks(model.TypePriceCounter, parseStockCounter); err != nil {
		return
	}
	fmt.Println("Crawl Dividends")
	if err = c.UpdateDividends(); err != nil {
		return
	}
	fmt.Println("Crawl Revenues")
	if err = c.UpdateRevenues(); err != nil {
		return
	}
	fmt.Println("End crawling")
	return nil
}

func (c Crawler) updateStocks(t model.Type, parse func([]string) model.Stock) (err error) {
	if !c.isExpire(t, "") {
		return nil
	}
	fmt.Println("Expire, crawl for new")
	file, err := download(filename(c.Config, t))
	if err != nil {
		return
	}

	if same, hash := c.isSame(file, t, ""); !same {
		reader := csv.NewReader(strings.NewReader(file))
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
		c.updateRecord(t, "", hash)
	} else {
		fmt.Printf("Data is same %s, do not thing\n", hash)
	}
	return nil
}

func filename(config config.Config, t model.Type) string {
	switch t {
	case model.TypePriceListed:
		return config.URL.Listed
	case model.TypePriceCounter:
		return config.URL.Counter
	}
	return ""
}

func download(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", ErrStockNotFound
	}
	if resp.StatusCode != 200 {
		return "", ErrStockNotFound
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return string(data), err
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
