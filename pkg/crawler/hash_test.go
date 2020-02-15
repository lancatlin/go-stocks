package crawler

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func getData(t *testing.T) []byte {
	page, err := download("https://matomo.wancat.cc")
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(page)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestHash(t *testing.T) {
	a := hashString(getData(t))
	time.Sleep(time.Second * 1)
	b := hashString(getData(t))
	fmt.Printf("%x\n%x\n", a, b)
	assert.Equal(t, a, b)
}
