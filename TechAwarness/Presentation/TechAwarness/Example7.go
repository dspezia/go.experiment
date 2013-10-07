package main

import (
	"fmt"
	"time"
)

// B1 OMIT
func main() {

	for i := 0; i < 4; i++ {
		go func(n int) {
			for j := 0; j < 3; j++ {
				fmt.Println("This is", n)
				time.Sleep(time.Second)
			}
			fmt.Println("Completed", n)
		}(i)
	}

	time.Sleep(5 * time.Second)
}

// E1 OMIT
