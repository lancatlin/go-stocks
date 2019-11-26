package model

type Dividend struct {
	StockID       string `gorm:"primary_key;type:varchar(20);auto_increment:false"`
	Year          int    `gorm:"primary_key;auto_increment:false"`
	MoneyDividend float64
	StockDividend float64
}
