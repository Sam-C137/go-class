# Interfaces & Methods in Detail

-> An interface variable is nil until initialized

It really has two parts:
* a value or pointer of some type
* a pointer to type information so the correct actual method can be identified

```text
var r io.Reader // nil until initialized
var b *bytes.Buffer // ditto

r = b // r is no longer nil but it has a nil pointer to a buffer
```

## Interfaces in practice

1. Let consumers define interfaces (what minimal behavior do they require?)
2. Re-use standard interfaces wherever possible (if interfaces from the standard library are used, our code becomes more compatible)
3. Keep the interface declarations small ("The bigger the interface, the weaker the abstraction")
4. Compose one-method interfaces into larger interfaces (if needed) eg: `io.Writter io.Reader -> io.WriteReader`
5. Avoid coupling interfaces to particular types/implementations
6. Accept interfaces, but return concrete types (let the consumer of the return type decide how to use it)

## Empty interfaces

The `interface{}` type has no methods, so is satisfied by anything. Empty interfaces are commonly used; they're how the formatted I/O
routines can print any type
```text
func fmt.Printf(f string, args ...interface{})
```

_Reflection_ is needed to determine what the concrete type is.