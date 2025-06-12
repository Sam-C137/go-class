package svl

type node struct {
	v int
	t *node
}

func insert(i int, h *node) *node {
	tail := &node{i, nil}

	if h != nil {
		h.t = tail
	}
	return tail
}

func makeList(n int) *node {
	var h, t *node
	h = insert(0, h)
	t = insert(1, h)

	for i := 2; i < n; i++ {
		t = insert(i, t)
	}
	return h
}

func sumList(h *node) (i int) {
	for n := h; n != nil; n = n.t {
		i += n.v
	}
	return i
}

func makeSlice(n int) []int {
	s := make([]int, n)

	for i := 0; i < n; i++ {
		s[i] = i
	}
	return s
}

func sumSlice(l []int) (i int) {
	for _, n := range l {
		i += n
	}
	return i
}
