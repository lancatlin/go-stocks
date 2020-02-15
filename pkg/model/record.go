package model

import "time"

type Type uint

const (
	TypePrice Type = iota
	TypeDividend
	TypeRevenue
)

type Record struct {
	ID        uint
	Type      Type
	StockID   string
	UpdatedAt time.Time
	ExpireAt  time.Time
}
