package crawler

import (
	"testing"

	"github.com/lancatlin/go-stocks/pkg/config"
)

func TestCrawlDividend(t *testing.T) {
	c := config.Config{}
	c.URL.Dividend = "https://tw.stock.yahoo.com/quote/%s/dividend"
	crawler := New(c)
	dividends := crawler.crawlDividend("2884")
	if len(dividends) == 0 {
		t.Fatal("Result length is 0:", dividends)
	}
}
