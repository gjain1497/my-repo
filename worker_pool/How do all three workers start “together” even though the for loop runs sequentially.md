How do all three workers start “together” even though the for loop runs sequentially?


The loop that starts workers
for w := 1; w <= 3; w++ {
    go worker(w, jobs, results, &wg)
}

Step 1: Loop runs sequentially

First iteration (w=1): go worker(1, ...)

Important: go immediately creates a new goroutine for worker1

Main goroutine continues to the next iteration immediately

Second iteration (w=2): go worker(2, ...)

Another goroutine created for worker2

Main continues immediately

Third iteration (w=3): go worker(3, ...)

Another goroutine created for worker3

Main finishes the loop

✅ So all 3 goroutines exist now.

Step 2: Why it looks like they “start together”

Each worker goroutine immediately executes its function (for job := range jobs) independently.

Even though main created them one by one, once created, the Go scheduler runs them concurrently.

Each worker hits <-jobs almost instantly, and blocks waiting for a job.

Think of it like: you call three people on the phone one by one, but as soon as they pick up, they’re all waiting on the same line simultaneously.

So in runtime diagrams, it’s correct to write:

Worker1 → blocks on <-jobs
Worker2 → blocks on <-jobs
Worker3 → blocks on <-jobs


even though the creation was sequential, the blocking happens concurrently.

Step 3: The key

Sequential creation ≠ sequential execution

go launches a new lightweight thread of execution

Scheduler decides when to run each, but since they all immediately block on <-jobs, it looks like they are “all ready at the same time.”

💡 Simple analogy:

Main is a manager calling 3 workers one by one to sit at a desk.

As soon as each worker sits, they put their headphones on and wait for tasks.

By the time the manager finishes calling all three, all 3 are sitting and waiting simultaneously — exactly what <-jobs does.

If you want, I can draw a tiny step-by-step timeline showing sequential creation vs concurrent blocking — it makes this click visually.