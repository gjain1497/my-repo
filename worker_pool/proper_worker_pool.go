package main

import (
	"fmt"
	"sync"
	"time"
)

type Task interface {
	Execute() (interface{}, error)
}

type SquareTask struct {
	Number int
}

func (t SquareTask) Execute() (interface{}, error) {
	// Simulate some work
	time.Sleep(100 * time.Millisecond)
	// Simulate random failures
	if t.Number%7 == 0 {
		return nil, fmt.Errorf("unlucky number %d", t.Number)
	}

	// Calculate square
	result := t.Number * t.Number
	return result, nil
}

type Job struct {
	ID       int
	Task     Task
	Priority int // optional (for later advanced features)
	Attempts int // how many times job was tried
}

type Result struct {
	JobID  int
	Err    error
	Output interface{}
}

type WorkerPool struct {
	NoOfWorkers    int
	JobsChannel    chan Job
	ResultsChannel chan Result
	wg             sync.WaitGroup
}

func NewWorkerPool(numWorkers int, jobQueuesize int) *WorkerPool {
	return &WorkerPool{
		NoOfWorkers:    numWorkers,
		JobsChannel:    make(chan Job, jobQueuesize),
		ResultsChannel: make(chan Result, jobQueuesize),
	}
}

func (w *WorkerPool) Start() {
	for i := 1; i <= w.NoOfWorkers; i++ {
		go w.worker(i) // Each worker runs in its own goroutine
	}
}
func (w *WorkerPool) worker(workerID int) {
	// Loop forever, processing jobs
	for job := range w.JobsChannel { // Blocks until job arrives
		//Execute the task
		output, err := job.Task.Execute()

		//Always send result(even if error occured, so
		//that it does not block other workers)
		w.ResultsChannel <- Result{
			JobID:  job.ID,
			Err:    err,
			Output: output,
		}
		w.wg.Done()

		// Optional: log
		if err != nil {
			fmt.Printf("Worker %d: Job %d failed: %v\n", workerID, job.ID, err)
		} else {
			fmt.Printf("Worker %d: Job %d completed\n", workerID, job.ID)
		}
	}
}

// Submit job to queue
func (w *WorkerPool) SubmitJob(job Job) error {
	//just send job  to channel
	select {
	case w.JobsChannel <- job:
		w.wg.Add(1)
		return nil
	default:
		return fmt.Errorf("job queue is full or pool is shut down")
	}
}

func (w *WorkerPool) Wait() {
	w.wg.Wait()
}

// Shutdown gracefully
func (wp *WorkerPool) Shutdown() {
	close(wp.JobsChannel)
}

func main() {
	//this will actually send all jobs to queue right? or maybe some diff process
	fmt.Println("=== Worker Pool Demo ===\n")

	// 1. Create pool with 3 workers and queue size 10
	pool := NewWorkerPool(10, 100)

	// 2. Start a goroutine to collect results
	go func() {
		for result := range pool.ResultsChannel {
			if result.Err != nil {
				fmt.Printf("❌ Job %d failed: %v\n", result.JobID, result.Err)
			} else {
				fmt.Printf("✅ Job %d result: %v\n", result.JobID, result.Output)
			}
		}
	}()
	// 3. Start the worker pool
	pool.Start()
	fmt.Println("Worker pool started with 3 workers\n")

	// 4. Submit 10 jobs
	fmt.Println("Submitting 10 jobs...")
	for i := 1; i <= 100; i++ {
		job := Job{
			ID:   i,
			Task: SquareTask{Number: i},
		}

		err := pool.SubmitJob(job)
		if err != nil {
			fmt.Printf("Failed to submit job %d: %v\n", i, err)
		}
	}
	fmt.Println("All jobs submitted!\n")

	pool.Wait()

	// 6. Shutdown the pool
	fmt.Println("\nShutting down pool...")
	pool.Shutdown()

	// 7. Give time for final results to print
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n=== Demo Complete ===")

}
