# Interfaces in HTTP

```text
type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}

type HandlerFunc func(w ResponseWriter, r *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}

// handler matches type HandlerFunc and so interface Handler
// the HTTP framework can call ServeHTTP on it

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world! from %s\n", r.URL.Path[1:])
}
```