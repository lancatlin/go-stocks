package web

import (
	"github.com/jinzhu/gorm"
	"github.com/lancatlin/go-stocks/pkg/model"
)

type RYG struct {
	model.Stock
	model.Revenue
	Returns []float64
}

func (r RYG) IsNil() bool {
	return r.Stock.ID == ""
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

	var revenue model.Revenue
	err := h.Where("stock_id = ?", id).Order("time desc").First(&revenue).Error
	if gorm.IsRecordNotFoundError(err) {
		h.AddRevenue(id)
		return h.RYG(id)
	}

	ryg := RYG{
		Stock:   stock,
		Revenue: revenue,
		Returns: make([]float64, 3),
	}
	for i, y := range []int{1, 5, 10} {
		ryg.Returns[i] = stock.ReturnOnInvestment(y)
	}
	return ryg
}
