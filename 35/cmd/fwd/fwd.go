package fwd

import "math/rand"

type forwarder interface {
	forward(string) int
}

type thing1 struct {
	t forwarder
}

type thing2 struct {
	t forwarder
}

type thing3 struct {
}

const defaultChars = "01234567 . . .mnopqrstuvwxyz"

func randString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

//go:noinline
func length(s string) int {
	return len(s)
}

func (t1 *thing1) forward(s string) int {
	return t1.t.forward(s)
}

func (t2 *thing2) forward(s string) int {
	return t2.t.forward(s)
}

func (t3 *thing3) forward(s string) int {
	return len(s)
}
