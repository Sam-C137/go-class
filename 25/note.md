# Context

→ Context was added in go as a way of providing a common method for cancellation 
of work in progress.

→ The `Context` package offers a common method to cancel requests.
- explicit cancellation
- implicit cancellation based on a timeout or deadline

→ A context may also carry request-specific values, such as a trace ID.
Many network or database requests, for example, take a context for cancellation.

→ Contexts form an immutable tree structure (go-routine safe; changes to a context
do not affect its ancestors). Cancellations or timeout applies to a current context and its
subtree. Ditto for value. 

→ A subtree may be created with a shorter timeout but not longer.

## Cancellation and timeouts

→ Context offers two controls:

- a channel that closes when the cancellation occurs
- an error that's readable once the channel closes

The error tells you whether the request was cancelled or timed out.
We often use the channel from `Done()` in a select block.

→ The context value should always be the first parameter of a function.

```text
// First runs a set of queries and returns the result from
// the first to respond, cancelling the others.
func First(ctx context.Context, urls []string) (*Result, error) {
    c := make(chanResult, len(urls))  // buffered to avoid orphans
	ctx, cancel := context.WithCancel(ctx)
	
	defer cancel()  // cancel the other queries when we are done
	
	search := func(url string) {
		c <- runQuery(ctx, url)
    }
	
	...
}
```

## Values

→ A context is also a way to pass around values. We can attach some values to a context; say
at the beginning of a query, and they are available in the various places we process the query below.

→ Context values should be data specific to a request, such as:
- a trace ID or start time (for latency calculation)
- security or authorization data

**Avoid** using context to carry optional parameters.

Use package-specific, private context key type (not string) to avoid collisions

```go
package main

import (
	"context"
	"log"
	"net/http"
)

type contextKey int

const TraceKey contextKey = 1

// AddTrace is a middleware to insert a traceId into a request
func AddTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if traceId := r.Header.Get("X-Cloud-Trace-Context"); traceId != "" {
			ctx = context.WithValue(ctx, TraceKey, traceId)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ContextLog makes a log with traceId as prefix
func ContextLog(ctx context.Context, f string, args ...any) {
	// -- reflection

	traceId, ok := ctx.Value(TraceKey).(string)

	if ok && traceId != "" {
		f = traceId + ": " + f
	}

	log.Printf(f, args...)
}
```