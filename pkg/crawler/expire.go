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
		expire = time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, time.Local)
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
		expire = time.Date(now.Year(), time.June, 1, 0, 0, 0, 0, time.Local)
		if now.After(expire) {
			expire = expire.AddDate(1, 0, 0)
		}

	case model.TypeRevenue:
		expire = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local)
	}
	return
}
