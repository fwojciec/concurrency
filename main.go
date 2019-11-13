package main

import (
	"fmt"
	"sync"
)

// GEN FUNCTION
// ------------
// a function like that can be used to generate workload tasks, which can be
// in turn taken up and completed by individual goroutines working in parallel.

// so in the case of our application this could be a function sending product
// hrefs on a chan string.

// or also, stepping back even further, a function that generates category
// links which will be only in turn used to generate product links.
func genV1(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// we can simplify this function by using a buffered channel.
// since the channel is buffered it can take all values without waiting or
// blocking. We can close the channel immediately after sending.
func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	for _, n := range nums {
		out <- n
	}
	close(out)
	return out
}

// this function is basically a "worker" function. it takes in a raw material
// of some sort, does work on it, and returns a value.

// in case of a scraper this could be a function that parses a page and returns
// individual data items.

// there can be multiple "worker" function taking messages from a single
// channel -- i.e. constituting a "fan out".
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			// "work"
			out <- n * n
		}
		close(out)
	}()
	return out
}

// this function is basically the "fan in" part of the arrangement. It takes in
// multiple channels as argument and returns values on a single channel.
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// start an output goroutine for each input channel in cs. output
	// copies values from c to out until c is closed, then calls wg.Done
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a gouroutine to close out once all the output gouroutines are
	// done. This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	in := gen(2, 3)

	c1 := sq(in)
	c2 := sq(in)

	for n := range merge(c1, c2) {
		fmt.Println(n)
	}
}
