package crawler

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lancatlin/go-stocks/pkg/model"
)

func (c Crawler) isExpire(t model.Type, id string) bool {
	var last model.Record
	err := c.Where("type = ? and stock_id = ? and expire_at > ?", t, id, time.Now()).First(&last).Error
	if gorm.IsRecordNotFoundError(err) {
		return true
	} else if err != nil {
		panic(err)
	}
	return false
}

func (c Crawler) updateRecord(t model.Type, id string) {
	record := model.Record{
		Type:      t,
		StockID:   id,
		UpdatedAt: time.Now(),
	}
	record.ExpireAt = expire(t)
	c.Save(&record)
}
