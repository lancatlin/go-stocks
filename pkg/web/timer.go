package web

import (
	"fmt"
	"time"
)

func (h Handler) UpdatePricesRegularly() {
	var last time.Time
	update := make(chan bool)
	keep := make(chan bool)
	reset := make(chan time.Time)
	go h.after(update, keep)
	for {
		select {
		case <-update:
			go func() {
				fmt.Println("receive event")
				if err := h.UpdateInfo(); err != nil {
					fmt.Println(err)
				}
				fmt.Println("\n\n\nupdate at", last, "\n\n\n")
				reset <- time.Now()
				keep <- true
			}()
		case <-h.ask:
			h.ans <- last
		case last = <-reset:
			fmt.Println("last is", last)
		}
	}
}

func (h Handler) after(callback, keep chan bool) {
	callback <- true
	<-keep
	for {
		select {
		case <-time.After(h.Update):
			callback <- true
			<-keep
		}
	}
}
