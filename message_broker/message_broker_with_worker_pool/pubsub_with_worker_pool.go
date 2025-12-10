package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Topic string
	Data  interface{}
}

type Topic struct {
	name        string
	subscribers []*Subscriber
	mu          sync.RWMutex
}

type PubSub struct {
	topics map[string]*Topic //[topic_name -> topic_object]
	mu     sync.RWMutex
}
type Task interface {
	Execute() (interface{}, error)
}

type ProcessMessageTask struct {
	msg     Message
	handler func(Message)
}

func (t *ProcessMessageTask) Execute() (interface{}, error) {
	t.handler(t.msg)
	return nil, nil
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

func NewPubSub() *PubSub {
	return &PubSub{
		topics: make(map[string]*Topic),
	}
}

func (ps *PubSub) CreateTopic(topicName string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	_, ok := ps.topics[topicName]
	if ok {
		return errors.New("Topic already exists")
	}
	topic := &Topic{
		name:        topicName,
		subscribers: []*Subscriber{},
	}
	ps.topics[topicName] = topic
	return nil
}

// Benefits:
// Don't hold PubSub lock while modifying topic
// Better concurrency (other threads can access other topics)
// Use RLock when only reading (multiple readers allowed
func (ps *PubSub) Subscribe(topicName string, subscriber *Subscriber) error {
	//Lock Pubsub to read topics map
	ps.mu.RLock() //Use RLock (only reading)
	topic, ok := ps.topics[topicName]
	ps.mu.RUnlock() //Unlock early

	if !ok {
		return errors.New("Topic does not exist")
	}

	//Lock Topic to modify subscribers
	topic.mu.Lock() //Lock the topic
	topic.subscribers = append(topic.subscribers, subscriber)
	topic.mu.Unlock() //Unlock the topic

	//Start subscriber
	subscriber.start() //Start listening
	return nil
}

func (ps *PubSub) Unsubscribe(topicName string, subscriberId string) error {
	// Remove subscriber from topic

	//Lock Pubsub to read topics map
	ps.mu.RLock() //Use RLock (only reading)
	topic, ok := ps.topics[topicName]
	ps.mu.RUnlock() //Unlock early

	if !ok {
		return errors.New("Topic does not exist")
	}

	//Lock Topic to modify subscribers
	topic.mu.Lock() //Lock the topic
	var updatedSubscribers []*Subscriber
	var removedSubscriber *Subscriber
	for _, subscriber := range topic.subscribers {
		if subscriber.id != subscriberId {
			updatedSubscribers = append(updatedSubscribers, subscriber)
		} else {
			removedSubscriber = subscriber
		}
	}
	topic.subscribers = updatedSubscribers
	topic.mu.Unlock() //Unlock the topic

	//close the removed subscriber's channel
	if removedSubscriber != nil {
		removedSubscriber.Close()
	} else {
		return errors.New("Susbcriber not found")
	}

	return nil
}

func (ps *PubSub) Publish(topicName string, data interface{}) error {
	//publish to all subscribers of this topic
	//getTopic from topic name
	ps.mu.RLock()
	topic, ok := ps.topics[topicName]
	ps.mu.RUnlock()

	if !ok {
		return errors.New("Topic does not exist")
	}

	topic.mu.RLock()
	//publish msg to all subscribers
	subscribers := topic.subscribers
	topic.mu.RUnlock()

	for _, subscriber := range subscribers {
		subscriber.channel <- Message{Topic: topicName, Data: data}
	}
	return nil
}

type Subscriber struct {
	id           string
	channel      chan Message
	handler      func(Message)
	WorkerPool   *WorkerPool
	jobIDCounter int
	mu           sync.Mutex
}

func NewSubscriber(id string, handler func(Message)) *Subscriber {
	return &Subscriber{
		id:         id,
		channel:    make(chan Message, 10),
		handler:    handler,
		WorkerPool: NewWorkerPool(5, 50),
	}
}

func (s *Subscriber) generateJobID() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobIDCounter++
	return s.jobIDCounter
}

func (s *Subscriber) start() {
	//start single goroutine that listens on the channel
	// go func() {
	// 	for msg := range s.channel {
	// 		s.handler(msg)
	// 	}
	// }()

	//Start worker pool, as in IMS we started when initalising it
	//because think of it like there was only one susbcriber
	//i.e IMS here. cmiiw (CLAUDE!)

	// IMS = 1 InventoryManager = 1 WorkerPool (started once)
	// PubSub = N Subscribers = N WorkerPools (each started once per subscriber)
	s.WorkerPool.Start()

	//Now instead of starting the handler directly we will
	//submit job to the specific task
	go func() {
		for msg := range s.channel {
			job := Job{
				ID: s.generateJobID(),
				Task: &ProcessMessageTask{
					msg:     msg,
					handler: s.handler,
				},
			}
			s.WorkerPool.SubmitJob(job)
		}
	}()
}

func (s *Subscriber) Close() {
	//stop the channel so nothing gets read from this
	close(s.channel)
	s.WorkerPool.Shutdown()
}

func main() {
	pubsub := NewPubSub()
	pubsub.CreateTopic("gaming")

	// Subscriber with SLOW handler (simulates heavy processing)
	sub1 := NewSubscriber("analytics", func(msg Message) {
		fmt.Printf("[Analytics] Started processing: %v\n", msg.Data)
		time.Sleep(2 * time.Second) // ⭐ Simulate 2 second processing
		fmt.Printf("[Analytics] Finished processing: %v\n", msg.Data)
	})

	pubsub.Subscribe("gaming", sub1)

	// Publish 10 messages rapidly
	fmt.Println("--- Publishing 10 messages ---")
	start := time.Now()

	for i := 1; i <= 10; i++ {
		pubsub.Publish("gaming", fmt.Sprintf("Message %d", i))
	}

	publishTime := time.Since(start)
	fmt.Printf("\nPublishing took: %v (fast!)\n", publishTime)

	// Wait for all workers to finish
	fmt.Println("\nWaiting for workers to process...")
	sub1.WorkerPool.Wait()

	totalTime := time.Since(start)
	fmt.Printf("\n✅ Total time: %v\n", totalTime)
	fmt.Printf("Expected: ~4 seconds (10 messages / 5 workers * 2 sec)\n")
	fmt.Printf("Without worker pool: ~20 seconds (10 messages * 2 sec)\n")
}
