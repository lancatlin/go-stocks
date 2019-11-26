package model

type Dividend struct {
	Stock         Stock  `gorm:"foreignkey:StockID"`
	StockID       string `gorm:"primary_key"`
	Year          int    `gorm:"primary_key"`
	MoneyDividend float64
	StockDividend float64
}