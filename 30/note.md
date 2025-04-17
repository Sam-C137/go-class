# Concurrency problems

1. Race conditions, where unprotected read & writes overlap
    * must be some data that is written to
    * could be a read-modify-write operation
    * and two goroutines can do it at the same time

2. Deadlock, when no goroutine can make progress
    * goroutines could all be blocked on empty channels
    * goroutines could all be blocked waiting on a mutex
    * GC could be prevented from running (busy loop)
   
   Go detects some deadlocks automatically; with `-race` it can find `some` data races

3. Goroutine leaks
    * goroutine hangs on an empty or blocked channel
    * not deadlock: other goroutines make progress
    * often found by looking at `ppfrof` output

    When you start a goroutine, always know how/when it will end

4. Channel errors
    * trying to send on a closed channel
    * trying to send or receive on a nil channel
    * closing a nil channel
    * closing a channel twice

5. Other errors
    * closure capture
    * misuse of `Mutex`
    * misuse of `WaitGroup`
    * misuse of `select`
   
    A good taxonomy of Go concurrency errors may be found in [this](https://cseweb.ucsd.edu/~yiying/GoStudy-ASPLOS19.pdf) paper.

### Data race

1.
```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

var nextID = 0

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>You got %v</h1>", nextID)
	//unsafe data race
	nextID++
}

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
```

### Deadlock
1.
```go
package main

import "fmt"

func main() {
	ch := make(chan bool)

	go func(ok bool) {
		fmt.Println("START")
		if ok {
			ch <- ok
        }
	}(false)
	
	<- ch
	fmt.Println("DONE")
}
```
Solution: Avoid impossible situations

2.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var m sync.Mutex
	done := make(chan bool)
	fmt.Println("START")

	go func() {
		m.Lock()
	}()

	go func() {
		time.Sleep(1)
		m.Lock()
		defer m.Unlock()
		fmt.Println("SIGNAL")
		done <- true
	}()
	
	<-done
	fmt.Println("DONE")
}
```
Solution: For every locked mutex make sure there is a corresponding unlock

3. Dining philosopher's problem

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var m1, m2 sync.Mutex
	done := make(chan bool)
	fmt.Println("START")

	go func() {
		m1.Lock()
		defer m1.Unlock()
		time.Sleep(1)
		m2.Lock()
		defer m2.Unlock()
		
		fmt.Println("SIGNAL")
		done <- true
	}()

	go func() {
		m2.Lock()
		defer m2.Unlock()
		time.Sleep(1)
		m1.Lock()
		defer m1.Unlock()

		fmt.Println("SIGNAL")
		done <- true
	}()
	
	<-done
	fmt.Println("DONE")
	<-done
	fmt.Println("DONE")
}
```
Solution: If you have to acquire multiple mutexes always make sure they are locked
and unlocked in the correct order

4. Goroutine leak

```txt
func finishReq(timeout time.Duration) *obj {
   ch := make(chan obj)
   
   go func() {
	     ...      // work that takes too long
	   ch <- fn() // blocking send
   }()

   select {
   case rslt := <- ch: 
	   return rslt
   case <- time.After(timeout):
	   return nil
   }
}
```
Solution: Buffer the channel so that after timeout if we happen to send to the channel
we can safely ignore the unused value

5. Incorrect use of waitgroup: always call `Add` before `go` or `Wait`

```txt
func walkDir(dir string, pairs chan<- pair, ...) {
   wg.Add(1) // BIG MISTAKE
   defer wg.Done()

   visit := func(p string, fi os.FileInfo) {
	   if fi.Mode().isDir() && p != dir {
		   go walkDir(p, pairs, wg. limits)
       ...
   }
}

err := walkDir(dir, paths, wg)
wg.Wait()
```
Solution: Always add to the wait group before you start the unit of work, i.e.,
```txt
     wg.Add(1)  
	 go walkDir(p, pairs, wg. limits)
```

6. Closure capture: a goroutine closure shouldn't capture a mutating variable

```go
package main

import "fmt"

func main() {
   for i := 0; i < 10; i++ { // WRONG
      go func() {
         fmt.Println(i)
      }()
   }
}
```
Solution: pass the variable's value as a parameter
```go
package main

import "fmt"

func main() {
   for i := 0; i < 10; i++ { // RIGHT
      go func(i int) {
         fmt.Println(i)
      }(i)
   }
}
```

### Select problems

`select` can be challenging and lead to mistakes:
* `default` is always active
* a `nil` channel is always ignored
* a full channel (for send) is skipped over
* a "done" channel is just another channel
* available channels are selected at random

1. Skipping a full channel to default and losing a message
```txt
for {
   x := socket.Read()

   select {
   case output <- x:
       ...
   default:
      return
   }
}
```
The code was written assuming we'd skip output only if it was set to nil
we also skip if output is full, and lose this and future messages

2. Reading a "done" channel and aborting when input is backed up by another 
channel—that input is lost
```txt
for {
   select {
   case x := <- input:
       ...
   case <- done:
      return
   }
}
```
There's no guarantee we read all of `input` before reading `done`
Solution: use done only for an error abort; close input on `EOF`
```txt
for {
   select {
   case x, ok := <- input:
       if !ok {
          return
       }
       ...
}
```

## Some thoughts

Four considerations when using concurrency:
1. Don't start a goroutine without knowing how it will stop
2. Acquire locks/semaphores as late as possible; release them in the reverse order
3. Don't wait for non-parallel work that you could do yourself
   ```text
      func do() int {
        ch := make(chan int)
        go func() { ch <- 1 }()
        return <- ch
      }
   ```
4. Simplify! Review! Test!