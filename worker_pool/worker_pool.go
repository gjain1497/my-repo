// // workerpool.go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// heavyCompute simulates a slow computation
func heavyCompute(n int) int {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(50))) // simulate work
	return n * n
}

func main() {
	// Example input (try 100000 to test big loads)
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Job queue and result queue
	jobs := make(chan int, len(numbers))
	results := make(chan int, len(numbers))

	// Start worker pool with 3 workers

	// What happens when Go sees go worker(...)

	// go keyword is encountered in the main goroutine.

	// Go runtime immediately schedules a new goroutine:

	// It sets up a new execution context (stack, registers, etc.) for worker(...).

	// That new goroutine starts running independently, concurrently with the main goroutine.

	// Main goroutine does NOT wait for the worker to run.

	// It continues immediately to the next statement in the loop (next iteration).

	// ✅ So yes — the main goroutine does not enter the function body itself.

	// Instead, it just tells the runtime: “Hey, run this function in a new goroutine.”

	// Control returns immediately to main.

	// for w in 1..3:
	// go worker(w)  <-- main says "new worker, run yourself"
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
	// main immediately moves to next iteration

	// Worker goroutines are like little robots: once called, they do their thing independently.

	// Main goroutine doesn’t care what the worker does; it just keeps looping.

	// Main: see go worker(1) → schedule W1 → move to next iteration
	// Main: see go worker(2) → schedule W2 → move to next iteration
	// Main: see go worker(3) → schedule W3 → done with loop

	// Meanwhile (concurrently):
	// W1 starts and blocks on <-jobs
	// W2 starts and blocks on <-jobs
	// W3 starts and blocks on <-jobs

	//g1<- job2 g2<-job1 g3<-job3

	for w := 1; w <= 3; w++ {
		go func(workerID int, jobs <-chan int, results chan<- int) {
			for n := range jobs {
				output := heavyCompute(n)
				fmt.Printf("Worker %d processed job %d → %d\n", workerID, n, output)
				results <- output
			}
		}(w, jobs, results)
	}
	// You asked:

	// "But we start workers one by one in the for loop. Isn’t that sequential? Shouldn’t workers run one after another instead of concurrently?"

	// Here’s the truth:

	// Yes, the loop in main is sequential. It creates worker goroutines one by one.

	// But the moment you call go worker(...), that worker is scheduled to run in parallel.

	// So after the loop, you don’t have “sequential workers”, you have 3 active concurrent workers, all waiting for jobs.

	// 💡 Think of it like: You hire workers one by one, but once hired, they all work at the same time.

	//after this loop
	//3 go routines got launched

	// Send all jobs to queue
	for _, n := range numbers {
		jobs <- n
	}
	close(jobs) // no more jobs

	// Collect all results
	for i := 0; i < len(numbers); i++ {
		res := <-results
		fmt.Println("Result:", res)
	}
}
