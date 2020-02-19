package web

import (
	"fmt"
	"time"
)

func (h Handler) UpdatePricesRegularly() {
	h.update()
	for {
		select {
		case <-time.After(time.Duration(h.Update) * time.Minute):
			h.update()
		}
	}
}

func (h Handler) update() {
	fmt.Println("receive event")
	if err := h.UpdateInfo(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("\n\n\nupdate at", formatTime(time.Now()), "\n\n\n")
}
