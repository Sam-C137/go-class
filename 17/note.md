# Object-oriented programming

→ For many people, the essential elements of object-oriented programming have been:
* abstraction
* encapsulation
* polymorphism
* inheritance

Go's approach to OOP is similar but different

## Abstraction

→ Decoupling behavior from the implementation details. The unix file system API is a greate example
of effective abstraction.
Roughly five basic functions hide all the messy details:
* open
* close
* read
* write
* ioctl

Many different operating system things can be treated like files

## Encapsulation

→ Hiding implementation details from misuse. It's hard to maintain an abstraction if the details are exposed:
* the internals may be manipulated in ways contrary to the concept behind the abstraction
* users of the abstraction may come to depend on the internal details—but those might change

Encapsulation usually means controlling the visibility of name s ("private" variables)

## Polymorphism

→ Polymorphism literally means "many shapes"—multiple types behind a single interface.
Three main types are recognized:
* ad-hoc: typically found in function/operator overloading
* parametric: commonly known as "generic programming" 
* subtype: subclasses substituting for superclasses

→ "Protocol-oriented" programming uses explicit interface types, now supported in many popular languages (an ad-hoc method).
In this case, behavior is completely separate from implementation which is good for abstraction 

## Inheritance

Inheritance has conflicting meanings:
* substitution (subtype) polymorphism
* structural sharing of implementation details

In theory, inheritance should always imply subtyping: the subclass should be "kind of" the superclass.
See the [Liskov substitution principle](https://reflectoring.io/lsp-explained/)

Theories about substitution can be pretty messy

→ Why would inheritance be bad?
It injects a dependence on the superclass into the subclass:
* what if the superclass changes behavior?
* what if the abstract concept is leaky?

Not having inheritance means better encapsulation and isolation

---------------------------------------------------------

→ Go offers four main supports for OOP
* encapsulation using the package for visibility control
* abstraction & polymorphism using interface types
* enhanced composition to provide structure sharing

Go does not offer inheritance
