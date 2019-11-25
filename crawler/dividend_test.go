package main

import (
	"testing"
)

func TestCrawDividend(t *testing.T) {
	id := "2884"
	t.Log(CrawlDividend(id))
}
