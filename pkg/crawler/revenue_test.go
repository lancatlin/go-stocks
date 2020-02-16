package crawler

import (
	"testing"

	"github.com/lancatlin/go-stocks/pkg/config"
)

func TestCrawlRevenue(t *testing.T) {
	id := "6112"
	conf := config.New()
	c := New(conf)
	revenue := c.crawlRevenue(id)
	t.Log(revenue)
}
