package model

import "time"

type Revenue struct {
	StockID      string    `gorm:"primary_key"`
	Time         time.Time `gorm:"primary_key"`
	MonthRevenue float64
	YearRevenue  float64
}
