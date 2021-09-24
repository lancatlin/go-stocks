package crawler

import (
	"fmt"
	"testing"

	"github.com/lancatlin/go-stocks/pkg/config"
)

func TestCrawlRevenue(t *testing.T) {
	c := config.Config{}
	c.URL.Revenue = "https://tw.stock.yahoo.com/quote/%s/revenue"
	crawler := New(c)
	revenue, err := crawler.crawlRevenue("2884")
	if err != nil {
		t.Fatal("Result length is 0:", err)
	}
	fmt.Println(revenue)
}
