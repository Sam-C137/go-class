# Odds & Ends
→ There are no real enumerated types in Go

→ You can make an almost-enum type using a named type and constants:

```go
package main

type shoe int

const (
	tennis shoe = iota
	dress
	sandal
	clog
)
```

iota starts at 0 in each const block and increments once on each line;
here 0, 1, 2, ..., the `shoe` type allows use to separate this enumerated type
from a normal integer

## Iota expressions

→ Traditional flags are easy:

```go
package main

type Flags uint

const (
	FlagUp Flags = 1 << iota // is up
	FlagBroadcast            // supports broadcast access
	FlagLoopback             // is a loopback interface
	FlagPointToPoint         // is a point-to-point link
	FlagMulticast            // supports multicast access
)
```

These flags take on the values in a power-of-two sequence: 0x01, 0x02, 0x04, etc.
That makes them easy to combine, e.g. `FlagUp` | `FlagLoopback`

→ Go also supports more complex `iota` expressions:

```go
package main

type ByteSize int64

const (
	_ = iota                        // ignore first value
	KiB ByteSize = 1 << (10 * iota) // 2 ^ 10
	MiB                             // 2 ^ 20 (1 << 10*2)
	GiB
	TiB
	PiB
	EiB
)
```

$2^{10}$ = 1024

So `EiB` is set to $2^{60}$ (1 << 10*6) = $1152921504606846976$ since it's the 6th value in the sequence after the ignored first value.

## Variable argument lists

→ What if we don't know how many parameters a function needs?

```go
package main

import "fmt"

func main() {
	fmt.Printf("%#v\n", make(map[string]int))
	fmt.Printf("%s: %s\n", "biege", "4")
}
```
All the formatted printing code uses variable argument lists

→ Variable argument functions are implemented:
```go
package main

import "fmt"

func sum(nums ...int) int {
    var total int
	for _, num := range nums {
		total += num
    }
	return total
}

func main() {
	fmt.Println(sum())
	fmt.Println(sum(1))
	fmt.Println(sum(1, 2, 3, 4))
	fmt.Println(sum([]int{1, 2, 3, 4}...))
}
```
Only the last parameter in a parameter list can be variable.

## Sized integers

→ Sometimes we need to handle low-level protocols (TCP/IP, etc.)

```go
package main

type TCPFields struct {
	SrcPort uint16
	DstPort uint16
	SeqNum uint32
	AckNum uint32
	DataOffset uint8
	Flags uint8
	WindowSize uint16
	Checksum uint16
	UrgentPtr uint16
}
```
So we need to work with integers that have a particular size and/or are unsigned

## Bitwise operators

```go
package main

import "fmt"

func main() {
	a, b := uint16(0xfff), uint16(281)

	fmt.Printf("%016b %#04[1]x\n", a)
	fmt.Printf("%016b %#04[1]x\n", a &^ 0b1111) // AND NOT
	fmt.Printf("%016b %#04[1]x\n", a & 0b1111)
	fmt.Println()
	fmt.Printf("%016b %#04[1]x\n", b) // ^ on a single operand in go acts like a bitwise NOT aka `~` in other languages
	fmt.Printf("%016b %#04[1]x\n", ^b)
	fmt.Printf("%016b %#04[1]x\n", b | 0b1111)
	fmt.Printf("%016b %#04[1]x\n", b ^ 0b1111)
}
```

We can combine the TCP declaration and an enumerated type:

```go
package main

// Flags that may be set in a TCP segment.
const (
	TCPFlagFin = 1 << iota
	TCPFlagSyn
	TCPFlagRst
	TCPFlagPsh
	TCPFlagAck
	TCPFlagUrg
)

// true if both flags are set
synAck := tcpHeader.Flags & (TCPFlagSyn|TCPFlagAck) == (TCPFlagSyn|TCPFlagAck)
```
Checking for bit flags this way is pretty common in low-level code

```go
package main

import "fmt"

func main() {
	a, b, c := uint16(1024), uint16(255), uint16(0xff00)

	fmt.Printf("%016b %#04[1]x\n", a)
	fmt.Printf("%016b %#04[1]x\n", a << 3) 
	fmt.Printf("%016b %#04[1]x\n", a << 13)
	fmt.Println()
	fmt.Printf("%016b %#04[1]x\n", b) 
	fmt.Printf("%016b %#04[1]x\n", b << 2)
	fmt.Printf("%016b %#04[1]x\n", b >> 2)
	fmt.Printf("%016b %#04[1]x\n", c >> 2)
}
```

## Signed short integers

```go
package main

import "fmt"

type short int16
type ushort uint16

func main() {
    a := short(-128)
	b, c := -a, 1/-1
	d, e := a+1, a-1
	
	fmt.Printf("%4d %#04x\n", a, ushort(a))
	fmt.Printf("%4d %#04x\n", b, ushort(b))
	fmt.Printf("%4d %#04x\n", c, ushort(c))
	fmt.Printf("%4d %#04x\n", d, ushort(d))
	fmt.Printf("%4d %#04x\n", e, ushort(e))
}
```

## `goto` considered harmful

→ Every once in a while, `goto` is simply easier to understand

```txt
readFormat:
	err := binary.Read(buf, binary.BigEndian, &header.format)
	
	if err != nil {
        return &header, nil, HeaderReadFailed.from(pos, err)
	}
	
	if header.format == junkID {
	    ... //find size & consume WAVE junk header
	    goto readFormat
	}
	
	if header.format != fmtID {
	    return &header, nil, InvalidChunkType
	}
```