package crawler

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lancatlin/go-stocks/pkg/model"
)

func (c Crawler) isExpire(t model.Type, id string) bool {
	var last model.Record
	err := c.Where(model.Record{
		Type:    t,
		StockID: id,
	}).Where("expire_at > ?", time.Now()).First(&last).Error
	if gorm.IsRecordNotFoundError(err) {
		return true
	} else if err != nil {
		panic(err)
	}
	return false
}

func (c Crawler) isSame(data interface{}, t model.Type, id string) (bool, string) {
	hash := hashString(data)
	var last model.Record
	err := c.Where(model.Record{
		Type:    t,
		StockID: id,
	}).Order("updated_at desc").First(&last).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, hash
	}
	return hash == last.Hash, hash
}

func (c Crawler) updateRecord(t model.Type, id, hash string) {
	record := model.Record{
		Type:      t,
		StockID:   id,
		Hash:      hash,
		UpdatedAt: time.Now(),
	}
	record.ExpireAt = expire(t)
	c.Save(&record)
}
