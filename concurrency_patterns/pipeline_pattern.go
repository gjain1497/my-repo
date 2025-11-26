package main

import (
	"fmt"
	"sync"
	"time"
)

func Generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// ⭐ Worker Pool style: 'go' outside
func Square(id int, in <-chan int, out chan<- int) {
	for n := range in {
		fmt.Printf("Worker %d processing: %d\n", id, n)
		time.Sleep(100 * time.Millisecond)
		out <- n * n
	}
	close(out)
}

func FanIn(workers []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range workers {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for result := range c {
				out <- result
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	start := time.Now()

	input := Generator(1, 2, 3, 4, 5)

	numWorkers := 3

	var results []<-chan int

	for i := 0; i < numWorkers; i++ {
		out := make(chan int)
		results = append(results, out)
		go Square(i, input, out) // ⭐ go outside just like worker pools

		//the point is we call it using go now
		//so we cannot get anything as a return type
		//from this function.

		//means we cannot write
		//ch := go Square(1, input)

		//so thats why there is no option
		//but to pass channel as paramter in the function
	}

	results2 := FanIn(results)
	for result := range results2 {
		fmt.Println("Result:", result)
	}

	elapsed := time.Since(start)
	fmt.Printf("\nTotal time: %v\n", elapsed)
}
