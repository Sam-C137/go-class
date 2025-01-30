# Networking with HTTP

→ The Go standard library has many packages for making web servers including:
* client and server sockets
* route multiplexing
* HTTP and HTML, including HTML templates
* JSON and other data formats
* cryptographic security
* SQL database access
* compression utilities
* image generation

→ There are also lots of 3rd-party packages with improvements

## Go HTTP design

→ An HTTP handler function is an instance of an interface

```
import (
    "net/http"
)


type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}

// The HTTP framework can call a method on a function type

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello world! from %s\n", r.URL.Path[1:])
}
```