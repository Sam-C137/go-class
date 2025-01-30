# Methods and Interfaces

## Why have methods?

→ An interface specifies absract behavior in terms of methods that a concrete object must provide

```go
package fmt

type Stringer interface {
	String() string
}
```

→ We define a method set on a concrete type and if the method set matches properly, then the
type can implement the interface.

## What are methods?
→ A method is a special type of function syntax. It has a reciever parameter before the function
name parameter.

```go
package main

import (
	"strconv"
	"strings"
)

type IntSlice []int

func (is IntSlice) String() string {
	var strs []string

	for _, v := range is {
		strs = append(strs, strconv.Itoa(v))
	}

	return "[" + strings.Join(strs, ";") + "]"
}
```
Now `IntSlice` implicitly implements `Stringer`

## Why interfaces

→ Without interfaces, we'd have to write may functions for many concrete types, possibly coupled
to them

```text
func OutputToFile(f *File, ...) { ... }
func OutputToBuffer(b *Buffer, ...) { ... }
func OutputToSocker(s *Socket, ...) { ... }
```

Better—we want to define our function in terms of abstract behavior

```text
type Writer interface {
	Write([]byte) (int, error)
}

func OutputTo(w io.Writer, ...)  { ... }
```

## Receivers
→ A method may take a pointer or value receiver, but not both

```go
package main

type Point struct {
	X, Y float64
}

func (p Point) Offset(x, y float64) Point {
    return Point{p.X + x, p.Y + y}
}

func (p *Point) Move(x, y float64) {
    p.X += x
	p.Y += y
}
```
Taking a pointer allows the method to change the receiver (original object)