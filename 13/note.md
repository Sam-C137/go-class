# Heap allocation

→ Go would prefer to allocate on the stack, but sometimes can't:
* a function returns a pointer to a local object
* a local object is captured in a function closure
* a pointer to a local object is sent via a channel
* any object is assigned into an interface
* any object whose size is variable at runtime (slices)

→ The use of new has nothing to do with this

→ Build with the flag `gcflags -m=2` t o see the escape analysis

## Slice safety

→ Anytime a function mutates a slice that is passed in, we must return a copy.
this is because if the original slice gets reallocated, all modifications we performed
on the slice goes away

```go
package main

func update(things []int) []int {
	things = append(things, 1)
	return things
}
``` 