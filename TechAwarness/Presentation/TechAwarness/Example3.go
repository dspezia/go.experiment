package main

import (
	"fmt"
)

func main() {
	// BEGIN OMIT
	var m = map[string]int{"foo": 0, "bar": 1}
	m["toto"] = 1
	delete(m, "bar")
	fmt.Println("Length: ", len(m), "of", m)
	x, b := m["do not exist"]
	fmt.Println("x:", x, "b:", b)
	// END OMIT
}
