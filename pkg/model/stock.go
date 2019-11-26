package model

type Stock struct {
	ID        string `gorm:"primary_key;type:varchar(20)"`
	Name      string
	Price     float64
	Dividends []Dividend
}
