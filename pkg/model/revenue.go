package model

import "time"

type Revenue struct {
	RevenueID    uint `gorm:"primary_key"`
	StockID      string
	Time         time.Time
	MonthRevenue float64
	YearRevenue  float64
}

func (r Revenue) Month() string {
	switch r.Time.Month() {
	case time.January:
		return "一月"
	case time.February:
		return "二月"
	case time.March:
		return "三月"
	case time.April:
		return "四月"
	case time.May:
		return "五月"
	case time.June:
		return "六月"
	case time.July:
		return "七月"
	case time.August:
		return "八月"
	case time.September:
		return "九月"
	case time.October:
		return "十月"
	case time.November:
		return "十一月"
	case time.December:
		return "十二月"
	default:
		return ""
	}
}
