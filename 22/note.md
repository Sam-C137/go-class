# Concurrency

→ Some definitions of concurrency

- "Execution happens in some non-deterministic order."
- "Undefined out-of-order execution"
- "Non-sequential execution"
- "Parts of a program execute out-of-order or in partial order"

Partial order

![](.\assets\partial-order.png)

- part 1 happens before parts 2 or 3
- both 2 and 3 complete before part 4
- the parts of 2 and 3 are ordered amon themselves

→ Subroutines and coroutines

Subroutines are subordinate, while coroutines are co-equal

![](.\assets\subroutine.png)

So let's try a new definition of concurrency:

**Parts of the program may be executed independently in some non-deterministic (partial) order**

## Parallelism

Parts of a program are executed independently at the same time. You can have concurrency with a single-core
processor (think interrupt handling in the operating system)

Parallelism can happen only on a multicore processor

Concurrency doesn't make the program faster, parallelism does

## Concurrency vs Parallelism

Concurrency is about dealing with things happening out-of-order, Parallelism is about things actually happening
at the same time

A single program won't have parallelism without concurrency We need concurrency to allow parts of the program to
execute independently. And that's where the fun begins...

## Race condition
"System behavior depends on the (non-deterministic) sequence or timing of parts of the program executing independently,
where some possible behaviors (orders of execution) produce invalid results"

To prevent a read-modify-write operation from causing out of sync error, we need to make them atomic; that means bundling
the parts together actively to prevent them from interleaving. Atom here, meaning non-divisible.

Solutions making sure operations produce a consistent state to any shared data
- don't share anything
- make the shared things read-only
- allow only one writer to the shared things
- make the read-modify-write operations atomic