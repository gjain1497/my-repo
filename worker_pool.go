package main

import (
	"time"
)

func heavyCompute(n int) int {
	time.Sleep(100 * time.Millisecond) // simulate heavy work
	return n * n
}

func main() {
	numJobs := 100_000 // large number to see memory spike
	numbers := make([]int, numJobs)
	for i := 0; i < numJobs; i++ {
		numbers[i] = i + 1
	}

	jobs := make(chan int, len(numbers)) //job queue

	results := make(chan int, len(numbers)) //result queue where workers will push results

	//2 Start workers
	for w := 1; w <= 3; w++ {
		go func(workerID int, jobs <- chan int, result <- chan int) {
			for n := range jobs { //keep reading until jobs channel is closed
				output := heavyCompute(n)
				results <- output
			}
		}(w, jobs, results)
	}
	//Each worker is just one go routine
	//It loops over jobs channel -> picks one job -> computes -> pushes results

	//3 Push jobs into the queue
	//send all jobs
	for _, n := range numbers{
		jobs <- n //instead of launching len(numbers) goroutines immediately, we push them into channel 
	}
	close(jobs) //signal workers no more jobs

	//Why? //instead of launching len(numbers) goroutines immediately, we push them into channel 
	//Workers will pull them one by one
	//closing channel is like saying that "all orders  placed no more new ones"


	//4 Collect results
	for i := 0; i < len(numbers); i++{
		res := <- results
	}


	



	// fmt.Println("Starting jobs without worker pool...")
	// start := time.Now()

	// var memStats runtime.MemStats
	// runtime.ReadMemStats(&memStats)
	// fmt.Printf("Memory before jobs: %v MB\n", bToMb(memStats.HeapAlloc))

	// for _, n := range numbers {
	// 	go func(num int) {
	// 		_ = heavyCompute(num)
	// 	}(n)
	// }

	// // Wait long enough (not ideal, just to observe memory)
	// time.Sleep(30 * time.Second)

	// runtime.ReadMemStats(&memStats)
	// fmt.Printf("Memory after jobs: %v MB\n", bToMb(memStats.HeapAlloc))
	// fmt.Println("All jobs done in", time.Since(start))
}

func bToMb(stats uint64) uint64 {
	return stats / 1024 / 1024
}
