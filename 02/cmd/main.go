package main

import (
	"02/hello"
	"fmt"
	"os"
)

// Go does not let the main function take in parameters for command line args
// instead we use os.Args

func main() {
	fmt.Printf(hello.Say(os.Args[1:]))
}
