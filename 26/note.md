# Channels In Detail

Channels as a communications tool and channels as a way for parts of your program
to sync with each other.

## Channel State

→ Channels block unless ready to read or write.

A channel is ready to write if:
- it has buffer space or
- at least one reader is ready to read

A channel is ready to read if:
- it has unread data in its buffer or
- at least one writer is ready to write or
- it is closed

→ Channels are unidirectional but have two ends (which can be passed separately as parameters)

- an end for writing and closing:
```txt
func get(url string, ch chan <- result)  { //write only end
    
}
```
- an end for reading
```txt
func collect(ch <-chan result) map[string]int  { //read only end
    
}
```

## Closed channels

→ Closing a channel causes it to return to the "zero" value.
We can receive a second value: is the channel closed?

```go
package main

import "fmt"

func main() {
	ch := make(chan int, 1)
	ch <- 1
	b, ok := <-ch
	fmt.Println(b, ok) // 1 true
	close(ch)
	c, ok := <- ch
	fmt.Println(c, ok) // 0, false
}
```

→ A channel can only be closed once (else it will panic).
One of the main issues in working with goroutines is ending them

- An unbuffered channel requires a reader and a writer (a writer blocked on a channel with no reader will "leak").
- Closing a channel is often a signal that work is done
- Only one goroutine can close a channel (not many)
- We may need a way to coordinate closing a channel or stopping goroutines (beyond the channel itself).

## `nil` channels

→ Reading or writing a channel that is `nil` always blocks*. But a `nil` channel in a select block is ignored. 

This can be a powerful tool:
- use a channel to get input
- suspend it by changing the channel variable to `nil`
- you can even unsuspend it again
- but really close the channel if there is no more input (EOF)

## Channel state reference

| State        | Receive         | Send          | Close                |
|--------------|-----------------|---------------|----------------------|
| Nil          | Block*          | Block*        | Panic                |
| Empty        | Block           | Write         | Close                |
| Partly full  | Read            | Write         | Readable until empty |
| Full         | Read            | Block         | Readable until empty |
| Closed       | Default Value** | Panic         | Panic                |
| Receive-only | Ok              | Compile Error | Compile Error        |
| Send-only    | Compile Error   | Ok            | Ok                   |

_*select ignores a nil channel since it would always block_
_**reading a closed channel returns (default value, !ok)_

## Buffering

→ Buffering allows the sender to send without waiting.

```go
package main

import "fmt"

func main() {
	// buffer that can hold 2 items
	messages := make(chan string, 2)

	// we can now send twice without getting blocked
	messages <- "buffered"
	messages <- "channel"

	// and receive both as usual
	fmt.Println(<-messages)
	fmt.Println(<-messages)
}
```
with a buffer size 1 or no buffer at all it will deadlock

→ Common uses of buffered channels:
- avoid goroutine leaks (from an abandoned channel)
- avoid rendezvous pauses (performance improvement)

→ Don't buffer until it's needed (buffering may hide a race condition!).

→ Some testing may be required to find the right number of slots

→ Special uses of buffered channels:
- counting semaphore pattern

### Counting semaphores
A counting semaphore limits work in progress (or occupancy). Once it is full,
only one unit of work can enter for each one that leaves. We model this with a 
buffered channel:

- attempt to send(write) before starting work
- the send will block if the buffer is full (occupancy is at max)
- receive (read) when the work is done to free up space in the buffer (this allows the next worker to start).