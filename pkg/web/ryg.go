package web

import (
	"github.com/jinzhu/gorm"
	"github.com/lancatlin/go-stocks/pkg/model"
)

type RYG struct {
	model.Stock
	Returns []float64
}

func (r RYG) IsNil() bool {
	return r.ID == ""
}

func (h Handler) RYG(id string) RYG {
	var stock model.Stock
	if err := h.Where("id = ?", id).Preload("Dividends", func(db *gorm.DB) *gorm.DB {
		return db.Order("dividends.year DESC")
	}).First(&stock).Error; gorm.IsRecordNotFoundError(err) {
		return RYG{}
	} else if err != nil {
		panic(err)
	}

	if len(stock.Dividends) == 0 {
		// if hadn't crawl yet
		h.UpdateDividend(id)
		return h.RYG(id)
	}

	ryg := RYG{
		Stock:   stock,
		Returns: make([]float64, 3),
	}
	for i, y := range []int{1, 5, 10} {
		ryg.Returns[i] = stock.ReturnOnInvestment(y)
	}
	return ryg
}
