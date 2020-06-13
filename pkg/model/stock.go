package model

type Stock struct {
	ID        string `gorm:"primary_key;type:varchar(20)"`
	Name      string
	Price     float64
	Dividends []Dividend `gorm:"PRELOAD:true"`
	Revenue   Revenue    `gorm:"PRELOAD:true"`
}

func (s Stock) ReturnOnInvestment(year int) float64 {
	var totalM, totalS float64
	for _, div := range s.Dividends[:min(year, len(s.Dividends))] {
		totalM += div.MoneyDividend
		totalS += div.StockDividend
	}
	avgM, avgS := totalM/float64(year), totalS/float64(year)/10
	newPrice := s.Price/(1+avgS) - avgM
	return (avgM + avgS*newPrice) / s.Price * 100
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
