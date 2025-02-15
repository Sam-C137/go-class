# Goroutines and channels

## Communicating sequential processes (CSP)

### Channels

→ A channel is a one way communications pipe.
- things go in one end, come out the other
- in the same order they went in
- until the channel is closed
- multiple readers and writers can share it safely

### Sequential process

→ Looking at a single independent part of the program, it
appears to be sequential

```text
for {
    read()
    process()
    write()
}
```
This is perfectly natural if we think of reading and writing files or network sockets

→ Putting the parts together we can use channels to communicate.

![](.\assets\img.png)

- each part is independent
- all they share are the channels between them
- the parts can run in parallel as the hardware allows

→ Why CSP is valuable:

→ Concurrency is always hard

→ CSP provides a model for thinking about it that makes it less hard
(take the program apart and make the pieces talk to each other)

→ "Go doesn't force devs to embrace the async ways of event driven programming.
It lets you write async code in a synchronous style."

## Goroutines

→ A goroutine is a unit of independent execution (coroutine)

→ It's easy to start a goroutine, put `go` in front of a function call

→ The trick is knowing how the goroutine will stop:

- you have a well-defined loop terminating condition, or
- you signal completion through a channel or context, or
- you let it run till the program stops

## Channels

→ A channel is like a one way socket or a unix pipe (except it allows for multiple readers and writers).
It is a method of synchronization as well as communication.
We know that a send (write) always happens before a receive (read)

→ It is also a vehicle for transferring ownership of data, so that only one goroutine at a time is writting
the data (avoid race conditions)

→ "Don't communicate by sharing memory, instead share memory by communicating" - Rob Pike

→
→
→
