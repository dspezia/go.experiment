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
	// END OMIT
}
