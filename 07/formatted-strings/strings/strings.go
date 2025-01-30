package strings

import "fmt"

func Strings() {
	s := "a string"
	b := []byte(s)

	// type
	fmt.Printf("%T\n", s)
	// value
	fmt.Printf("%v\n", s)
	// go formated value
	fmt.Printf("%#v\n", s)
	// quoted string
	fmt.Printf("%q\n", s)

	fmt.Printf("%v\n", string(b))
}
