package kafka

// Layer 1: KAFKA BROKER (Server)
// ============================================
// FILE: kafka-broker/main.go
// Run: go run kafka-broker/main.go
// Port: 9092
// ============================================

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

// ============================================
// MESSAGE & REQUEST STRUCTURES
// ============================================

type Message struct {
	Offset int64  // Position in partition
	Key    string // Used for routing
	Value  string // Actual data
}

type Request struct {
	Type       string // "FETCH", "PRODUCE", "SUBSCRIBE", "COMMIT"
	Topic      string
	Partition  int
	Offset     int64
	Key        string
	Value      string
	GroupID    string
	ConsumerID string
}

type Response struct {
	Type      string
	Messages  []Message
	Partition int
	Error     string
}

// ============================================
// BROKER - The Kafka Server
// ============================================

type Broker struct {
	topics         map[string]*Topic
	consumerGroups map[string]*ConsumerGroup // ← Track offsets ONLY
	mu             sync.RWMutex
}

type Topic struct {
	name       string
	partitions []*Partition
}

type Partition struct {
	id       int
	messages []Message // ← THE QUEUE (log)
	mu       sync.RWMutex
}

// Consumer Group - tracks committed offsets
// This is just BOOKKEEPING/TRACKING
// NOT the actual consumer!
type ConsumerGroup struct {
	groupID          string
	committedOffsets map[int]int64 // partition_id -> last committed offset
	mu               sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		topics:         make(map[string]*Topic),
		consumerGroups: make(map[string]*ConsumerGroup),
	}
}

// ============================================
// CREATE TOPIC
// ============================================

func (b *Broker) CreateTopic(name string, numPartitions int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	partitions := make([]*Partition, numPartitions)
	for i := 0; i < numPartitions; i++ {
		partitions[i] = &Partition{
			id:       i,
			messages: []Message{},
		}
	}

	b.topics[name] = &Topic{
		name:       name,
		partitions: partitions,
	}

	fmt.Printf("[Broker] Created topic '%s' with %d partitions\n", name, numPartitions)
}

// ============================================
// PRODUCE - Producer writes message
// ============================================

func (b *Broker) Produce(topic, key, value string) error {
	b.mu.RLock()
	t, exists := b.topics[topic]
	b.mu.RUnlock()

	if !exists {
		return fmt.Errorf("topic not found")
	}

	// Determine partition by hashing key
	partitionID := b.hashKey(key) % len(t.partitions)
	partition := t.partitions[partitionID]

	partition.mu.Lock()
	defer partition.mu.Unlock()

	// Create message
	offset := int64(len(partition.messages))
	msg := Message{
		Offset: offset,
		Key:    key,
		Value:  value,
	}

	// Append to partition log
	partition.messages = append(partition.messages, msg)

	fmt.Printf("[Broker] Stored message in %s-partition-%d at offset %d\n",
		topic, partitionID, offset)

	// ✅ Message is stored
	// ❌ Broker does NOT push to consumers
	// ❌ Broker doesn't even know who the consumers are!
	// ⏳ Message sits here until consumer pulls it

	return nil
}

// ============================================
// FETCH - Consumer pulls messages
// ============================================

func (b *Broker) Fetch(topic string, partition int, offset int64) ([]Message, error) {
	b.mu.RLock()
	t, exists := b.topics[topic]
	b.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("topic not found")
	}

	if partition >= len(t.partitions) {
		return nil, fmt.Errorf("partition not found")
	}

	p := t.partitions[partition]

	p.mu.RLock()
	defer p.mu.RUnlock()

	// Consumer asked: "Give me messages starting from offset X"
	start := int(offset)

	if start >= len(p.messages) {
		// No new messages available
		return []Message{}, nil
	}

	// Return up to 10 messages
	end := start + 10
	if end > len(p.messages) {
		end = len(p.messages)
	}

	messages := p.messages[start:end]

	fmt.Printf("[Broker] Returned %d messages from %s-partition-%d (offset %d to %d)\n",
		len(messages), topic, partition, start, end-1)

	return messages, nil
}

// ============================================
// GET COMMITTED OFFSET
// ============================================

func (b *Broker) GetCommittedOffset(groupID string, partition int) int64 {
	b.mu.RLock()
	group, exists := b.consumerGroups[groupID]
	b.mu.RUnlock()

	if !exists {
		return 0 // Start from beginning
	}

	group.mu.RLock()
	defer group.mu.RUnlock()

	offset, exists := group.committedOffsets[partition]
	if !exists {
		return 0 // Start from beginning
	}

	return offset
}

// ============================================
// COMMIT OFFSET - Consumer saves progress
// ============================================
func (b *Broker) CommitOffset(groupID string, partition int, offset int64) error {
	b.mu.Lock()
	group, exists := b.consumerGroups[groupID]
	if !exists {
		// Create new consumer group
		group = &ConsumerGroup{ // This is just BOOKKEEPING/TRACKING
			// NOT the actual consumer!
			groupID:          groupID,
			committedOffsets: make(map[int]int64),
		}
		b.consumerGroups[groupID] = group //b
	}
	b.mu.Unlock()

	group.mu.Lock()
	defer group.mu.Unlock()

	group.committedOffsets[partition] = offset

	fmt.Printf("[Broker] Group '%s' committed offset %d for partition %d\n",
		groupID, offset, partition)

	return nil
}

func (b *Broker) hashKey(key string) int {
	hash := 0
	for _, c := range key {
		hash = hash*31 + int(c)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}

// ============================================
// NETWORK SERVER
// ============================================

func (b *Broker) Start(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[Broker] Kafka listening on %s\n", port)

	// Create default topic
	b.CreateTopic("orders", 3)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		fmt.Printf("[Broker] Client connected from %s\n", conn.RemoteAddr())

		// Handle client requests
		go b.handleClient(conn)
	}
}

func (b *Broker) handleClient(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var req Request
		err := decoder.Decode(&req)
		if err != nil {
			fmt.Printf("[Broker] Client disconnected\n")
			return
		}

		var resp Response

		switch req.Type {
		case "FETCH":
			// Consumer pulling messages
			messages, err := b.Fetch(req.Topic, req.Partition, req.Offset)
			if err != nil {
				resp.Error = err.Error()
			} else {
				resp.Messages = messages
			}

		case "PRODUCE":
			// Producer writing message
			err := b.Produce(req.Topic, req.Key, req.Value)
			if err != nil {
				resp.Error = err.Error()
			}

		case "COMMIT":
			// Consumer committing offset
			err := b.CommitOffset(req.GroupID, req.Partition, req.Offset)
			if err != nil {
				resp.Error = err.Error()
			}

		case "GET_OFFSET":
			// Consumer asking for last committed offset
			offset := b.GetCommittedOffset(req.GroupID, req.Partition)
			resp.Messages = []Message{{Offset: offset}}
		}

		resp.Type = req.Type
		encoder.Encode(resp)
	}
}

func main() {
	broker := NewBroker()
	broker.Start(":9092")
}
