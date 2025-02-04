# Composition

→ Go mainly does composition with embedded structs. THe fields of an embedded struct are
promoted to the level of the embedding structure

```text
type Pair struct {
	Path string
	Hash string
}

type PairWithLength struct {
	Pair
	Length int
}

pl := PairWithLength{Pair{"/usr", "0xfdfe"}, 121}
fmt.Println(pl.Path, pl.Length) // not pl.x.Path
```

## Composition with pointer types

→ A struct can embed a pointer to another type; promotion of its fields and methods works
the same way

```text
type PairWithLength struct {
	Pair
	Length int
}

type Fizgig struct {
	*PairWithLength
	Broken bool
}

fg := Fizgig {
    &PairWithLength{Pair{"/usr", "0xfdfe"}, 121},
    false,
}

fmt.Println(fg) 
// Length of /usr is 121 with hash 0xfdfe
```

## Sortable interface

→ sort.Interface and sort.Sort are defined as

```go
package sort

type Interface interface {
	// Len is the number of elements in the collection
    Len() int
	
	// Less reports whether the element with 
	// index i should sort before the element with index j.
	Less(i, j int) bool
	
	// Swap swaps the elements with indexes i and j
	Swap(i, j int)
}

func Sort(data Interface)  {}
```

### Sortable built-ins

Slices of strings can be sorted using StringSlice

```text
	entries := []string{"charlie", "able", "dog", "baker"}
	// type cast entries to a type that has sort.Interface defined
	sort.Sort(sort.StringSlice(entries))
	fmt.Println(entries)
```

### Sorting example

Implement sort.Interface to make a type sortable

```go
package main

import (
	"fmt"
	"slices"
)

type Organ struct {
	Name   string
	Weight int
}

type Organs []Organ

func (s Organs) Len() int      { return len(s) }
func (s Organs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByName struct{ Organs }
type ByWeight struct{ Organs }

func (s ByName) Less(i, j int) bool {
	return s.Organs[i].Name < s.Organs[j].Name
}

func (s ByWeight) Less(i, j int) bool {
	return s.Organs[i].Weight < s.Organs[j].Weight
}

func main() {
	s := []Organ{{"brain", 1340}, {"liver", 1494}, {"spleen", 162}, {"pancreas", 131}, {"heart", 290}}
	slices.Sort(ByWeight{s})
	fmt.Println(s)
	slices.Sort(ByName{s})
	fmt.Println(s)
}
```

## Sorting in reverse

→ Use sort.Reverse which is defined as:

```text
package sort

type reverse struct {
	// This embedded Interface permits Reverse to use the methods of
	// another Interface implementation.
    Interface
}

// Less returns the opposite of the embedded implementation's Less method.
func (r reverse) Less(i, j int) bool {
    return r.Interface.Less(j, i)
}

// Reverse returns the reverse order for data.
func Reverse(data Interface) Interface  {
    return &reverse{data}
}
```