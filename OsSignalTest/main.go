package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	sc := make(chan os.Signal)
	signal.Notify(sc, os.Interrupt)

loop:
	for {
		select {
		case <-sc: // 여기서 ctrl +c 신호가 들어온다.
			fmt.Println("interrupt")
			break loop
		case <-time.After(2 * time.Second):
			fmt.Println("timeout", time.Now())
		}
	}

}
