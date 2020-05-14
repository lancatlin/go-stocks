package web

import (
	"github.com/lancatlin/go-stocks/pkg/model"
)

type RYG struct {
	model.Stock
	Returns []float64
	Error   string
}

func (h Handler) RYG(id string) RYG {
	stock, err := h.GetStock(id)
	if err != nil {
		return RYG{
			Error: err.Error(),
		}
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
