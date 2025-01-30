package hello

import "testing"

// To run all tests from the root, we run `go test ./...`

func TestHello(t *testing.T) {
	subtests := []struct {
		items  []string
		result string
	}{
		{
			result: "Hello, stranger!",
		},
		{
			items:  []string{"test"},
			result: "Hello, test!",
		},
		{
			items:  []string{"foo", "bar", "baz"},
			result: "Hello, foo, bar, baz!",
		},
	}

	for _, st := range subtests {
		if s := Say(st.items); s != st.result {
			t.Errorf("wanted %s (%v), got %s", st.result, st.items, s)
		}
	}
}
