# Select
→ The `select` statement is a powerful statement specifically for concurrency.
`select` is another control structure like `if`s and `for`s that allow us to work with channels
and goroutines, not other kinds of logic.

→ `select` works on channels and allows us to multiplex channels. `select` allows any ready alternative
to proceed among:
- a channel we can read from
- a channel we can write to
- a default action that's always ready

→ Most often, `select` runs in a loop so we keep trying. We can put a timeout or done channel into the 
select.

We can compose channels as synchronization primitives! Traditional primitives (mutex, condition variables)
can't be composed.

→ Taking this example

```go
package main

func main() {
	chans := []chan int{
		make(chan int),
		make(chan int),
	}
	for i := 0; i < 12; i++ {
		v1 := <- chans[0] 
		v2 := <- chans[1] 
	}
}
```
Reading from the channels like this means until we get a value from channel1, 
we never read from channel 2. This means if one channel sends faster than the other, we are going to have a
problem. Instead, this can be fixed with a select statement which allows us to listen to both channels at the same
time

```text
for i := 0; i < 12; i++ {
	select {
	case v1 := <-chans[0]:
		fmt.Println("received ", v1)
	case v2 := <-chans[1]:
		fmt.Println("received ", v2)
	}
}

```

## Default

→ In a select block, the default case is always ready and will be chosen if no other case is:

```go
package main

import "log"

var ch = make(chan []byte)

func sendOrDrop(data []byte) {
	select {
	case ch <- data:
		// sent ok; do nothing
	default:
		log.Printf("overflow: dropped %d bytes", len(data))
	}
}
```

Don't use default inside a loop, the select will be busy and waste cpu.