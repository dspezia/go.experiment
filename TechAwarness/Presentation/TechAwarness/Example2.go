package main

import (
	"fmt"
)

func main() {
	// BEGIN OMIT
	var t [4]int = [4]int{1, 2, 3, 4} // t is an array
	var s []int = t[1:3]              // s is a slice
	fmt.Println("s:", s, "len:", len(s), "cap:", cap(s))
	s = append(s, 10)
	fmt.Println("s:", s, "and t:", t, "!!!")
	// END OMIT
}
