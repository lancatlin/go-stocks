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
	"time"

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
	if err = c.importToDatabase(model.TypePriceListed, c.Config.URL.Listed, parseStockListed); err != nil {
		return
	}
	if err = c.importToDatabase(model.TypePriceCounter, c.Config.URL.Counter, parseStockCounter); err != nil {
		return
	}
	if err = c.UpdateDividends(); err != nil {
		return
	}
	if err = c.UpdateRevenues(); err != nil {
		return
	}
	fmt.Println("End crawling")
	return nil
}

func (c Crawler) isExpire(t model.Type) bool {
	var last model.Record
	err := c.Where("type = ? and expire_at > ?", t, time.Now()).First(&last).Error
	if gorm.IsRecordNotFoundError(err) {
		return true
	} else if err != nil {
		panic(err)
	}
	fmt.Println(last)
	return false
}

func (c Crawler) importToDatabase(t model.Type, filename string, parse func([]string) model.Stock) (err error) {
	if !c.isExpire(t) {
		return nil
	}
	file, err := download(filename)
	if err != nil {
		return
	}

	if same, hash := c.isDataSame(file, t); !same {
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
		c.updatePriceRecord(t, hash)
	} else {
		fmt.Println("Data is same", hash)
	}
	return nil
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

func (c Crawler) isDataSame(file string, t model.Type) (bool, string) {
	hash := hashString(file)
	var last model.Record
	err := c.Where("type = ?", t).Order("updated_at desc").First(&last).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, hash
	}
	return hash == last.Hash, hash
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

func (c Crawler) updatePriceRecord(t model.Type, hash string) {
	now := time.Now()
	record := model.Record{
		Type:      t,
		Hash:      hash,
		UpdatedAt: now,
	}
	expire := time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, time.Local)
	if now.Hour() > 14 {
		expire = expire.AddDate(0, 0, 1)
	}
	record.ExpireAt = expire
	c.Save(&record)
}
