# Error Handling

→ Errors in go are objects satisfying the `error` interface:

```txt
type error interface {
	 Error() string
}
```

→ Any concrete type with `Error()` can represent an error

```txt
type fizgig struct {}

func (e fizgig) Error() string{
    return "Your fizgig is bent"
}
```

## Building a custom error type

→ We're going to build out a custom error type

```go
package main

import (
	"fmt"
)

type errKind int

const (
	_ errKind = iota
	noHeader
	cantReadHeader
	invalidHdrType
	invalidChkLength
	invalidLength
)

// And we have some prototype errors we can return or customize
var (
	HeaderMissing      = WaveError{kind: noHeader}
	HeaderReadFailed   = WaveError{kind: cantReadHeader}
	InvalidHeaderType  = WaveError{kind: invalidHdrType}
	InvalidChunkLength = WaveError{kind: invalidChkLength}
	InvalidDataLength  = WaveError{kind: invalidLength}
)

type WaveError struct {
	kind  errKind
	value int
	err   error
}

func (e WaveError) Error() string {
	switch e.kind {
	case noHeader:
		return "no header (file too short?)"
	case cantReadHeader:
		return fmt.Sprintf("can't read header[%d]: %s", e.value, e.err)
	case invalidHdrType:
		return "invalid header type"
	case invalidChkLength:
		return fmt.Sprintf("invalid chunck length: %d", e.value)
	case invalidLength:
		return fmt.Sprintf("invalid data length: %d", e.value)
	default:
		return e.err.Error()
	}
}

// We have a couple of helper methods to generate errors

// with a particular value (e.g., header type)
func (e WaveError) val(val int) WaveError {
	temp := e // do a copy
	temp.value = val
	return temp
}

// `from` returns an error with a particular location and
// underlying error (e.g., from the standard library)
func (e WaveError) from(pos int, err error) WaveError {
	temp := e
	temp.value, temp.err = pos, err
	return temp
}
```
Here is an example of those errors in use

```txt
func DecodeHeader(b []byte) (*Header, []byte, error) {
	var err error
	var pos int

	header := Header{TotalLength: uint32(len(b))}
	buf := bytes.NewReader(b)
	if len(b) < headerSize {
		return &header, nil, HeaderMissing
	}

	if err := binary.Read(buf, binary.BigEndian, &header.riff); err != nil {
		return &header, nil, HeaderReadFailed.from(pos, err)
    }
}
```

## Wrapped errors

→ Starting `go 1.13` we can wrap one error in another

```txt
func (h HAL9009) OpenPodBayDoors() error {
    ...
    if h.err != nil {
        fmt.Errorf("I am sorry %s, I cannot %w", h.victim, h.err)
    }
    ...
} 
```
→ The easiest way to do that is to use the `%w` format verb with `fmt.Errorf()`

→ Wrapping errors gives us an error chain we can unravel
![](.\assets\img.png)

→ Custom error types may now unwrap their internal errors

```txt
type WaveError struct {
    value int
	err error
}

func (w WaveError) Unwrap() error {
    return w.err
}
```

## `errors.Is`

→ We can check whether an error has another error in its chain.
`errors.Is` compares with an error **variable** not a type

```txt
...
    if audio, err := DecodeWaveFile(fn); err != nil {
        if errors.Is(err, os.ErrPermission) {
            // let's report a security violation
        }
    }
...
```

→ We can provide the `Is` method for our custom error type (only useful if we export our error variables also)

```txt
type WaveError struct {
    kind errKind
    ...
}

func (w *WaveError) Is(t error) bool {
    e, ok := t.(*WaveError) // reflection again
    if !ok {
        return false
    }
    
    return e.errKind == w.errKind
}
```

## `errors.As`

→ We can get an error of an underlying type if it's in the chain. 
`errors.As` looks for an error **type** not a value.

```txt
...
    if audio, err := DecodeWaveFile(fn); err != nil {
        var e os.PathError // a struct
        
        if errors.As(err, &e) {
            // let's pass back the underlying file error
            return e
        }
    }
...
```

## A philosophy of error handling

→ When it comes to errors you may fall into one of these camps:
1. you hate constantly writing if/else blocks
2. you think writing if/else blocks makes things clearer
3. you don't care because you're too busy writing code

### Normal errors

→ Normal errors result from input or external conditions (for example a "file not found" error),
go handles this case by returning the error type.

```txt
// Not exactly os.Open but shows the basic logic
func Open(name string, flag int, perm FileMode) (*File, error) {
	r, e := syscall.Open(name, flag|syscall.O_CLOEXEC, syscallMode(perm))

	if e != nil {
        return nil, PathError{"open", name, e}
	}
	
	return NewFile(uintptr(r), name, kindOpenFile), nil
}
```

### Abnormal errors

→ Abnormal errors result from invalid program logic (for example a nil pointer).
For program logic errors, go does a panic.

```txt
func (d *digest) checkSum() [Size]byte {
    //finish writing the checksum
    ...
    if d.nx != 0 {
        panic("d.nx != 0") // panic if there's data left over
    }
    ...
}
```

→ When you program has a logic bug; if your server crashes it will get immediate attention
* logs are often noisy
* so proactive log searches for problems are rare

→ We want the evidence of failure as close as possible in time and space to the original 
defect in the code
* connect the crash to logs that explain the context
* traceback from the point closest to the broken logic

→ In a distributed system crash failures are the safest type to handle
* it's better to die than to be a zombie or babble or corrupt the db
* not crashing may lead to _Byzantine_ failures

## When should we panic?

→ Only when the error was caused by our own programming defect, e.g;
* we can't walk a data structure we built
* we have an off-by-one bug encoding bytes

In other words

"_panic should be used when our assumptions of our own programming design or logic are wrong_"

these cases may use an `assert` in other programming languages

→ A B-tree data structure satisfies several invariants:
1. every path from the root to a leaf has the same length
2. if a node has $n$ children, it contains $n - 1$ keys
3. every node (except the root) is at least half full
4. the root has at least two children if it is not a leaf
5. subnode keys fall between the keys of the parent node that lie on either side of the subnode pointer

if any of these is ever false, the B-tree methods should **panic!**

## Exception handling

→ Exception handling was popularized to allow "graceful degradation"
of safety-critical systems (e.g., Ada and flight control software).

→ Ironically, most safety-critical software are built without using exceptions!

→ Exception handling introduces invisible control paths through code so code with exceptions
is harder to analyze (automatically or by eye)

→ Officially go does not support exception handling as in other languages

→ Practically it does—in the form of `panic` and `recover`.

→ `panic` in a function will still cause deferred function calls to run, then it
will stop only if it finds a valid `recover` call in a `defer` as it unwinds the stack

## Panic and recover

→ `recover` from `panic` only works inside `defer`

```go
package main

import "fmt"

func abc() {
	panic("omg")
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			// what else can you do ?
			fmt.Println("recover: ", p)
		}
	}()
	
	abc()
}
```

## Define errors out of existence

→ Error (edge) cases are one of the primary sources of complexity. The best
way to deal with many errors is to make them impossible

→ Design your abstractions so most (or all) are safe
* reading from a nil map
* appending to a nil slice
* deleting a non-existent item from a map
* taking the length of an uninitialized string

Try to reduce edge cases that are hard to test or debug (or even think about)

## Proactively preventing problems

→ Every piece of data in your software should start life in a valid state

→ Every transformation should leave it in a valid state
* break large programs into small pieces you can understatnd
* hide information to reduce the chance of corruption
* avoid clever code and side effects
* avoid unsafe operations
* assert your invariants
* never ignore errors
* test, test, test

Never accept input from a user (or environment) without validation

## Error handling culture

→ "Go programmers think about the failure case first.
We solve the 'what if ...' case first. This leads to programs where failures are handled
at the point of writing, rather than the point they occur in production. The verbosity of
```txt
if err != nil {
    return err
}
```
is outweighed by the value of deliberately handling each failure condition at the point which
it occurs. Key to this is the cultural value of handling each and every error explicitly." — Dave Cheney
