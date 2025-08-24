package main

import (
	"errors"
	"fmt"
)

var (
	ErrQueueFull       = errors.New("queue is full")
	ErrQueueEmpty      = errors.New("queue is empty")
	ErrInvalidCapacity = errors.New("capacity must be positive")
)

// CircularQueue represents a circular queue data structure
type CircularQueue struct {
	items    []int
	front    int
	rear     int
	size     int
	capacity int
}

// NewCircularQueue creates a new CircularQueue with the given capacity
func NewCircularQueue(capacity int) (*CircularQueue, error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	return &CircularQueue{
		items:    make([]int, capacity),
		front:    0,
		rear:     -1, // Start at -1 so first enqueue sets rear to 0
		size:     0,
		capacity: capacity,
	}, nil
}

// Size returns the current number of elements in the queue
func (q *CircularQueue) Size() int {
	return q.size
}

// Capacity returns the maximum capacity of the queue
func (q *CircularQueue) Capacity() int {
	return q.capacity
}

// IsEmpty checks if the queue is empty
func (q *CircularQueue) IsEmpty() bool {
	return q.size == 0
}

// IsFull checks if the queue is full
func (q *CircularQueue) IsFull() bool {
	return q.size == q.capacity
}

func (q *CircularQueue) Metadata() {
	fmt.Printf("Front: %d, Rear: %d, Size: %d, Capacity: %d\n", q.front, q.rear, q.size, q.capacity)
}

// Front returns the front element without removing it (peek operation)
func (q *CircularQueue) Front() (int, error) {
	if q.IsEmpty() {
		return 0, ErrQueueEmpty
	}
	return q.items[q.front], nil
}

// Rear returns the rear element without removing it
func (q *CircularQueue) Rear() (int, error) {
	if q.IsEmpty() {
		return 0, ErrQueueEmpty
	}
	return q.items[q.rear], nil
}

func (q *CircularQueue) Enqueue(item int) error {
	if q.IsFull() {
		return ErrQueueFull
	}

	q.rear = (q.rear + 1) % q.capacity
	q.items[q.rear] = item
	q.size++
	return nil
}

func (q *CircularQueue) Dequeue() (int, error) {
	if q.IsEmpty() {
		return 0, ErrQueueEmpty
	}

	item := q.items[q.front]

	// update front pointer and size
	q.front = (q.front + 1) % q.capacity
	q.size--

	return item, nil
}

func (q *CircularQueue) Display() {
	fmt.Printf("Visual: ")
	for i := 0; i < q.capacity; i++ {
		marker := " "
		if i == q.front && i == q.rear && q.size > 0 {
			marker = "FR" // Front and Rear at same position
		} else if i == q.front && q.size > 0 {
			marker = "F"
		} else if i == q.rear && q.size > 0 {
			marker = "R"
		}

		if q.size > 0 {
			// Check if this index contains an active element
			isActive := false
			if q.front <= q.rear {
				isActive = i >= q.front && i <= q.rear
			} else {
				isActive = i >= q.front || i <= q.rear
			}

			if isActive {
				fmt.Printf("[%d%s] ", q.items[i], marker)
			} else {
				fmt.Printf("[ %s] ", marker)
			}
		} else {
			fmt.Printf("[ %s] ", marker)
		}
	}
	fmt.Println()
}

func main() {
	q, _ := NewCircularQueue(5)

	// Enqueue 3 items
	q.Enqueue(10)
	q.Enqueue(20)
	q.Enqueue(30)
	// Enqueue 2 more items
	q.Enqueue(40)
	q.Enqueue(50)

	q.Display()

	// Dequeue 2 items
	q.Dequeue()
	q.Display()

	q.Dequeue()
	q.Display()

	// Enqueue 1 item (should wrap around)
	q.Enqueue(60)
	q.Display()

	q.Enqueue(70)
	q.Display()

	err := q.Enqueue(80) // This should fail as the queue is full
	if err != nil {
		fmt.Println("After enqueueing 80:", err)
	}

	// Dequeue all
	for !q.IsEmpty() {
		item, err := q.Dequeue()
		if err != nil {
			fmt.Println("Error dequeuing:", err)
		} else {
			fmt.Println("Dequeued:", item)
		}
	}
	q.Display()

	_, err = q.Dequeue()
	if err != nil {
		fmt.Println("Error dequeuing:", err)
	}
	q.Display()
}
