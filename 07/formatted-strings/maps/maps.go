package maps

import "fmt"

func Maps() {
	m := map[string]int{"and": 1, "or": 2}

	// type
	fmt.Printf("%T\n", m)
	// value
	fmt.Printf("%v\n", m)
	// go formated value
	fmt.Printf("%#v\n", m)
}
