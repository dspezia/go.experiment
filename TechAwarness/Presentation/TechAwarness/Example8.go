package main

import (
	"fmt"
	"time"
)

// B1 OMIT
func main() {

	ch := make(chan string)

	go func() {
		for j := 0; j < 3; j++ {
			fmt.Println(j)
			time.Sleep(time.Second)
		}
		ch <- "hello"
	}()

	s := <-ch
	fmt.Println(s)
}

// E1 OMIT

func dummy() {

	normalTraffic := make(chan string)
	adminTraffic := make(chan string)

	// B2 OMIT
	select {
	case msg1 := <-normalTraffic:
		fmt.Println("Processing normal traffic", msg1)
	case msg2 := <-adminTraffic:
		fmt.Println("Processing admin command", msg2)
	case <-time.After(5 * time.Second):
		fmt.Println("Inactivity time out (5 seconds)")
	}
	// E2 OMIT
}
