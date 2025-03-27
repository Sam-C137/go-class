# Conventional Synchronization

→ These are primitives that are available in most programming languages, we will discuss
the `go` implementation.

→ Package `sync` for example:
- Mutex
- Once
- Pool
- RWMutex
- WaitGroup

→ Package `sync/atomic` for atomic scalar reads & writes

# Mutual exclusion

→ The problem we are trying to solve is mutual exclusion. What if multiple goroutines must read 
& write some data?

→ We must make sure only one of them can do so at any instant (in the so-called "critical section")

→ We accomplish this with some type of lock:
- acquire the lock before accessing the data
- any other goroutine will block waiting to get the lock
- release the lock when done

→ Problem: race condition

```go
package main

import (
	"fmt"
	"sync"
)

func do() int {
	var n int64
	var wg sync.WaitGroup

	for range 1000 {
		wg.Add(1)
		go func() {
			n++ // DATA RACE
			wg.Done()
		}()
	}

	wg.Wait()
	return int(n)
}

func main() {
	fmt.Println(do())
}
```
randomly prints numbers that may be less than a 1000

→ Solution: use a lock or mutex

```txt

func do() int {
	var n int64
	var wg sync.WaitGroup
	var mu sync.Mutex

	for range 1000 {
		wg.Add(1)
		go func() {
			mu.Lock()
			n++ // SAFE
			mu.Unlock()
			wg.Done()
		}()
	}

```

## Mutexes in action

```go
package main

import "sync"

type SafeMap struct {
	sync.Mutex     // not safe to copy
	m map[string]int
}

// so methods must take a pointer, not a value
func (s *SafeMap) Incr(key string) {
    s.Lock()
	defer s.Unlock()
	
	// only one goroutine can execute this
	// code at the same time, guaranteed
	s.m[key]++
}
```
Using `defer` is a good habit - avoid mistakes

→ It has a certain performance overhead. When we use a mutex, typically we are going to drop into the operating system, and it has the same
effect whether we are reading or writing, and sometimes it may cause a thread to get rescheduled.

## RWMutexes in action

→ Sometimes we need to prefer readers to (infrequent) writers

```go
package main

import (
	"sync"
	"time"
)

type InfoClient struct {
	mu        sync.RWMutex
	token     string
	tokenTime time.Time
	TTL time.Duration
}

func (i *InfoClient) CheckToken() (string, time.Duration) {
    i.mu.RLock()
	defer i.mu.RUnlock()
	
	return i.token, i.TTL - time.Since(i.tokenTime)
}
```

→ What makes the `RWMutex` different is that it prefers readers and it allows multiple readers to read at the same time. When locked for writing
all the readers have to be blocked, but when locked for reading, other readers are okay a writer would block till unlocked but other readers are okay
to read and that's very efficient.

```text

func (i *InfoClient) ReplaceToken(ctx context.Context) (string, error) {
    token, ttl, err := i.getAccessToken(ctx)
	if err != nil {
        return "", err
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	i.token = token
	i.tokenTime = time.Now()
	i.TTL = time.Duration(ttl) * time.Second
	
	return token, nil
}
```

## Atomic primitives

```go
package main

import (
	"sync"
	"sync/atomic"
)

func do() int {
	var n int64
	var wg sync.WaitGroup

	for range 1000 {
		wg.Add(1)
		go func() {
			atomic.AddInt64(&n, 1) // fixed
			wg.Done()
		}()
	}
	
	wg.Wait()
	return int(n)
}
```
This depends a lot on the underlying hardware to do an atomic add.

## Only-once execution

→ A `sync.Once` object allows us to ensure a function runs only once (only the first call to Do will call the function passed in)

```txt
var once sync.Once
var x *singleton

func initialize() {
    x = NewSingleton()
}

func handle(w http.ResponseWriter, r *http.Request) {
    once.Do(initialize)
	...
}
```
Checking `x == nil` in the handler is **unsafe!**

## Pool

→ A `Pool` provides for efficient & safe reuse of objects, but it's a container of `interface{}`

```go
package main

import (
	"bytes"
	"io"
	"sync"
)

var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func Log(w io.Writer, key, val string) {
    b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	//write to it
	w.Write(b.Bytes())
	bufPool.Put(b)
}
```
Pool is a device that allows us to not just manage memory, but to manage it safely in a concurrent program

## Wait, there's more!

→ Other primitives:
- Condition variable
- Map/`sync.Map` (safe container; uses `interface{}`)
- WaitGroup