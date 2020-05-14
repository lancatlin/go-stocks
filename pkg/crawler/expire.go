package crawler

import (
	"time"

	"github.com/lancatlin/go-stocks/pkg/model"
)

// expire define the logic of expire time
func expire(t model.Type) (expire time.Time) {
	now := time.Now()
	switch t {
	case model.TypePriceCounter, model.TypePriceListed:
		expire = time.Date(now.Year(), now.Month(), now.Day(), 14, 30, 0, 0, time.Local)
		if now.After(expire) {
			expire = expire.AddDate(0, 0, 1)
		}
		switch expire.Weekday() {
		case time.Saturday:
			expire = expire.AddDate(0, 0, 2)
		case time.Sunday:
			expire = expire.AddDate(0, 0, 1)
		}

	case model.TypeDividend:
		expire = now.AddDate(0, 0, 1)

	case model.TypeRevenue:
		expire = now.AddDate(0, 0, 1)
	}
	return
}
