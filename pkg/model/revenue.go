package model

import "time"

type Revenue struct {
	StockID      string
	Time         time.Time
	MonthRevenue float64
	YearRevenue  float64
}
