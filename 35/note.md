# Go benchmarks

→ Go has standard tools and conventions for running benchmarks. Benchmarks live in test files ending
with `_test.go`. You run benchmarks with `go test -bench`. Go only runs the `Benchmark{name}` functions.

→ Flags
* `-benchtime` - change benchmark time, default is `1s` the longer the time, the more accurate benchmarks
* `-benchmem` - benchmark memory and see how many allocations happen
* `-cpu` - 

→ Few other things to consider
* is the data/code available in cache?
* did you hit a garbage collection?
* did virtual memory have to page in/out?
* did branch prediction work the same way?
* did the compiler remove code via optimization? (are there side effects in the code?)
* are you running in parallel? how many cores?
* are those cores physical or virtual?
* are you sharing a core with anything else?
* what other processes are sharing the machine?

→ Code TOC:
fwd - Method forwarding
share - False sharing
svl - Slice vs List