package model

type Stock struct {
	ID           string `gorm:"primary_key;type:varchar(20)"`
	Name         string
	Price        float64
	MonthRevenue float64
	YearRevenue  float64
	Dividends    []Dividend `gorm:"PRELOAD:true"`
}

func (s Stock) ReturnOnInvestment(year int) float64 {
	var totalM, totalS float64
	for _, div := range s.Dividends[:min(year, len(s.Dividends))] {
		totalM += div.MoneyDividend
		totalS += div.StockDividend
	}
	avgM, avgS := totalM/float64(year), totalS/float64(year)
	return (avgM/s.Price + avgS/10) * 100
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
