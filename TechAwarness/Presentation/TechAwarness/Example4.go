package main

import (
	"fmt"
)

func foo() interface{} {
	x := 3
	return &x
}

func open(s string) (int, error) {
	return 0, nil
}

func readfile(f int) {
	return
}

func main() {

	// B1 OMIT
	if f, err := open("file"); err != nil {
		fmt.Println("Error:", err)
	} else {
		readfile(f)
	}
	// E1 OMIT

	// B2 OMIT
	s := []int{1, 2, 3, 4}
	for i := 0; i < len(s); i++ {
		fmt.Println(i, s[i])
	}
	for i, v := range s {
		fmt.Println(i, v)
	}
	// E2 OMIT

	x := 2

	// B3 OMIT
	// No break clause needed
	switch x {
	case 3, 4:
		fmt.Println("3 or 4")
	case 5, 6:
		fmt.Println("5 or 6")
	default:
		fmt.Println("Default")
	}

	// Arbitrary conditions can also be used
	switch {
	case x <= 100:
		fmt.Println("Small")
	case x > 100 && x < 1000:
		fmt.Println("Medium")
	default:
		fmt.Println("Large")
	}
	// E3 OMIT

	i := foo()

	// B4 OMIT
	switch i.(type) {
	case nil:
		fmt.Println("type is interface{}")
	case int, int64:
		fmt.Println("type is integer")
	case string:
		fmt.Println("type is a string")
	default:
		fmt.Println("I don't care")
	}
	// E4 OMIT

}

// B5 OMIT
func CalculateSomething(x int, s string) (int, bool) {
	return x + len(s), true
}

// E5 OMIT

// B6 OMIT
func NextInt() func() int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}

// E6 OMIT
