package crawler

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func getData(t *testing.T) string {
	page, err := download("https://matomo.wancat.cc")
	if err != nil {
		t.Fatal(err)
	}
	return page
}

func TestHash(t *testing.T) {
	a := hashString(getData(t))
	time.Sleep(time.Second * 1)
	b := hashString(getData(t))
	fmt.Printf("%x\n%x\n", a, b)
	assert.Equal(t, a, b)
}
