package crawler

import (
	"crypto/md5"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lancatlin/go-stocks/pkg/model"
)

func hashString(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func (c Crawler) isDivSame(id string, divs []model.Dividend) (bool, string) {
	hash := hashString(fmt.Sprint(divs))
	fmt.Println(hash)
	var last model.Record
	err := c.Where("type = ? and stock_id = ?", model.TypeDividend, id).Order("updated_at desc").First(&last).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, hash
	} else if err != nil {
		panic(err)
	}
	return hash == last.Hash, hash
}
