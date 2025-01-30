package numbers

import "fmt"

func Numbers() {
	a, b := 12, 345
	c, d := 1.2, 3.25

	// decimal
	fmt.Printf("%d %d\n", a, b)
	// hexadecimal. will print hex letters in lowercase if x is lowercase or vice versa
	fmt.Printf("%x %x\n", a, b)
	// creates a useful output more like we'd see in code eg: `0xc 0x159`
	fmt.Printf("%#x %#x\n", a, b)

	// floating point values, `%.(n)f` specifies number of decimal places
	fmt.Printf("%f %.2f\n", c, d)

	fmt.Println()

	// this says print this column in a width that is 6 characters
	fmt.Printf("|%6d|%6d|\n", a, b)
	// adding a 0 formats it with leading 0's
	fmt.Printf("|%06d|%06d|\n", a, b)
	// left justify
	fmt.Printf("|%-6d|%-6d|\n", a, b)
}
