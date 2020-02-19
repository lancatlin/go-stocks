package web

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lancatlin/go-stocks/pkg/model"
)

const (
	RED    = "#ff3c38"
	YELLOW = "#fde74c"
	GREEN  = "#81e979"
)

func getColor(n float64) string {
	switch {
	case n < 5:
		return RED
	case 5 <= n && n < 8:
		return YELLOW
	case 8 <= n:
		return GREEN
	}
	return ""
}

func percent(n float64) string {
	return fmt.Sprintf("%.2f", n)
}

func formatTime(t time.Time) string {
	zone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
	return t.In(zone).Format("2006-01-02 15:04:05")
}

func (h Handler) lastPriceUpdated() (listed, counter time.Time) {
	return h.lastRecord(model.TypePriceListed).UpdatedAt, h.lastRecord(model.TypePriceCounter).UpdatedAt
}

func (h Handler) lastRecord(t model.Type) (last model.Record) {
	err := h.Where("type = ?", t).Order("updated_at desc").First(&last).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		panic(err)
	}
	return
}
