package slices

import "fmt"

func Slices() {
	s := []int{1, 2, 3}
	a := [3]rune{'a', 'b', 'c'}

	// type
	fmt.Printf("%T\n", s)
	// value
	fmt.Printf("%v\n", s)
	// go formated value
	fmt.Printf("%#v\n", s)

	// type
	fmt.Printf("%T\n", a)
	// value
	fmt.Printf("%v\n", a)
	// go formated value
	fmt.Printf("%#v\n", a)
	// prints a string with quote marks
	fmt.Printf("%q\n", a)
}
