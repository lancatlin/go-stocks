package model

import "time"

type Revenue struct {
	RevenueID    uint `gorm:"primary_key"`
	StockID      string
	Time         time.Time
	MonthRevenue float64
	YearRevenue  float64
}
