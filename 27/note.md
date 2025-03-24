# Concurrent Sequential Processing

## What is the problem?

→ We have a folder containing ~50k files, and we want to find and remove
duplicates based on their content.

→ Use a secure hash because names/dates may differ

```text
fo88193 2
    /Users/foo/Dropbox/Emergency/FEMA_P-320_2014_508.pdf
    /Users/foo/Dropbox/Emergency/nps61-072915-01.pdf
```

## Normal approach

→ Checkout [normal](cmd/normal/main.go)

## Concurrent approach

→ The bulk of this work lies within reading all the individual files and hashing them.

→ We can use a fixed pool of workers who are all reading from a channel; this channel would
contain paths returned by the `path.Walk` function, then these goroutines can hash the individual
files concurrently and when their done, they feed an individual collector who's job is to manage the final hash table.

→ We can use a couple of approaches

1. [sequential](cmd/sequential/main.go)

2. [parallel directories](cmd/parallel-directories/main.go) ―
add a go routine for each directory in the tree •
this improves performance slightly; we're not waiting on paths to be identified.
3. [go routines galore](cmd/semaphore/main.go) ―
use a goroutine for every directory and file hash • what could go wrong? without some
controls, we'll run out of threads • solution; limit the number of active goroutines using a counting semaphore

### Evaluation
→ Using 32 workers was the best time. Increasing the limits buffer makes the time grow longer due to disk
contention

→ Amdahl's law: speedup is limited by the part (not) parallelized

$$
S = \frac{1}{1 - p + (p/s)}
$$

Here we've managed about $S = 6.25$ on $s = 8$ processors, or about $p = 96%$ parallel

![](.\assets\img.png)

## Conclusions

→ We don't need to limit goroutines

→ We need to limit contention for shared resources; the go runtime will limit the number of threads for the cpu normally through go.MAXPROCS,
and so if we are CPU bound, it's not really a problem, but when we started doing I/O bound work in this case the disk, then we need to limit contention
to that I/O resource because that's not scheduled, and we did that using a counting semaphore and that gave us the best performance.
