package model

import "time"

type Type uint

const (
	TypePriceListed Type = iota
	TypePriceCounter
	TypeDividend
	TypeRevenue
)

type Record struct {
	ID        uint
	Type      Type
	StockID   string
	Hash      string
	UpdatedAt time.Time
	ExpireAt  time.Time
}
