package model

type Stock struct {
	ID        string `gorm:"primary_key"`
	Name      string
	Price     float64
	Dividends []Dividend
}

