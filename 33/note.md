# Reflection

## Type assertion

→ `interface{}` says nothing since it has no methods. It's a generic thing, but sometimes
we need its "real" type. We can extract a specific type with a type assertion (aka downcasting).
This has the form `value.(T)` for some type `T`

```txt
    var w io.Writer = os.Stdout
	f := w.(*os.File) // success: f == os.Stdout
	c := w.(*bytes.Buffer) // panic: interface holds *os.File, not *bytes.Buffer
```
If we use the two-result version, we can avoid panic

```txt
    var w io.Writer = os.Stdout
	f, ok := w.(*os.File) // success: ok, f == os.Stdout
	b, ok := w.(*bytes.Buffer) // failure: !ok, b == nil
```

## Deep equality

→ We can use the `reflect` package in UIs to check equality

```txt
    want := struct {
		a: "a string",
		b: []int{
		1, 2, 3
	} // not comparable with ==
	}

	got ;= gotGetIt(...)
	
	if !reflect.DeepEqual(got, want) {
		t.Errorf("bad response: got=%#v, want=%#v", got, want)
    }
```
You can use `github.com/kylelemon/godebug/prettty` to show a deep diff

## Switching on type

→ We can also use type assertion in a `switch` statement (matching a type not a value)

```txt
func Println(args ...interface{}) {
    buff := make([]byte, 0, 80)
	for arg := range args {
		switch a := arg.(type) {
		case string: // concrete type
            buff = append(buff, a...)
		case Stringer: // interface
			buf = append(buff, a.String()...)
			...
		}
    }
}
```
Here the switch variable `a` has a specific type if the case has a single type

## Hard JSON

→ Not all JSON messages are well-behaved. What if some keys depend on others in the message?

```json
{
  "item": "album",
  "album": {"title":  "Dark side of the Moon"}
}
```
```json
{
  "item": "song",
  "song": {"title": "Bella Donna", "artist": "Stevie Nicks"}
}
```

### Custom JSON decoding
→ We'll make a wrapper and a custom decoder

```go
package main

type response struct {
    Item string `json:"item"`
	Album string
	Title string
	Artist string
}

type respWrapper struct {
	response
}

// We need respWrapper because it must have a separate unmarshal
// method from the response type
// check main.go for full implementation
```

## Testing JSON
→ We want to know if a known fragment of JSON is contained in a larger unknown piece
`{"id": "Z"}` in? `{"id": "Z", "part": "fizgig", "qty": 2}` 
All done with reflection from a generic map

```go
package main

func matchNum(key string, exp float64, data map[string]any) bool {
    if v, ok := data[key]; ok {
		if val, ok := v.(float64); ok && val == exp {
			return true
        }
    }
	
	return false
}
```
