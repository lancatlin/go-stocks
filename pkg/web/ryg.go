package web

import (
	"github.com/lancatlin/go-stocks/pkg/model"
)

type RYG struct {
	model.Stock
	Returns []float64
}

func (r RYG) IsNil() bool {
	return r.Stock.ID == ""
}

func (h Handler) RYG(id string) RYG {
	stock, err := h.GetStock(id)
	if err != nil {
		return RYG{}
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
