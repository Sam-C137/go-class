package hello

import (
	"strings"
)

func Say(names []string) string {
	if len(names) < 1 {
		names = []string{"stranger"}
	}

	return "Hello, " + strings.Join(names, ", ") + "!"
}
