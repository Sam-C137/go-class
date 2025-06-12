package share

import "sync"

const (
	nworker = 8
	buffer  = 1024
)

var wg sync.WaitGroup

func fill(n int, in chan<- int) {
	for i := 0; i < n; i++ {
		in <- i
	}

	close(in)
}

func run() (total int) {
	cnt := make([]uint64, nworker) // one cache line
	in := make([]chan int, nworker)

	for i := 0; i < nworker; i++ {
		in[i] = make(chan int, buffer)
		go fill(10000, in[i])
	}

	for i := 0; i < nworker; i++ {
		wg.Add(1)
		go count(&cnt[i], in[i])
	}

	wg.Wait()

	for _, v := range cnt {
		total += int(v)
	}

	return
}

func count(cnt *uint64, in <-chan int) {
	// false sharing
	for i := range in {
		*cnt += uint64(i)
	}

	wg.Done()
}

//func count(cnt *uint64, in <-chan int) {
//	var total int
//
//	for i := range in {
//		total += i
//	}
//
//	*cnt = uint64(total)
//	wg.Done()
//}
