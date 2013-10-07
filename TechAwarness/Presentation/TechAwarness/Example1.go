package main

import "fmt"

const PI = 3.14

func main() {

	fmt.Println("Hello, 世界")

	var sum int
	for i := 0; i < 10; i++ {
		sum += i
	}
	if sum > 20 {
		fmt.Println("sum=", sum)
	}
	for sum > 0 {
		fmt.Printf("%f\n", float64(sum)*PI)
		sum /= 2
	}
}
