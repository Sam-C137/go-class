# Standard I/O

→ Unix has the notion of three standard I/O streams

→ They're open by default in every program

→ Most modern programming languages have followed this convention:
* Standard input
* Standard output
* Standard error (output)

→ These are normally mapped to the console/terminal but can be redirected
``` bash
find . -name '*.go' | xargs grep -n "rintf" > print .txt
```

→ In go you can do a `fmt.Println()` and it goes by default to the standard
output, but we can redirect this to the standard error output by doing:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("printing a line to standard output")

	fmt.Fprintln(os.Stderr, "printing to error output")
}
```

→ Package `os` has functions to open or create files, list directories, etc. and
hosts the `File` type

→ Package `io` has utilities to read and write; `bufio` provides the buffered I/O 
scanners, etc.

→ Package `io/ioutil` has extra utilities such as reading an entire file to memory,
or writing it out all at once

→ Package `strconv` has utilities to convert to/from string representations