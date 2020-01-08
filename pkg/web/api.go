package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (h Handler) searchStock(c *gin.Context) {
	type Stock struct {
		ID   string
		Name string
	}
	var stocks []Stock
	if c.Query("q") == "" {
		c.Status(404)
		return
	}
	query := fmt.Sprintf("%%%s%%", c.Query("q"))
	err := h.Limit(10).Where("id LIKE ?", query).Or("name LIKE ?", query).Find(&stocks).Error
	if err != nil {
		panic(err)
	}
	c.JSON(200, stocks)
}
