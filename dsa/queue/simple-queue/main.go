package main

import "fmt"

type Queue struct {
	Items []int
}

func NewQueue() *Queue {
	return &Queue{
		Items: []int{},
	}
}

func (q *Queue) IsEmpty() bool {
	return len(q.Items) == 0
}

func (q *Queue) Size() int {
	return len(q.Items)
}

func (q *Queue) Display() {
	fmt.Printf("Queue: %v (front -> rear)\n", q.Items)
}

func (q *Queue) Enqueue(item int) {
	q.Items = append(q.Items, item)
}

func (q *Queue) Dequeue() (int, bool) {
	if q.IsEmpty() {
		return 0, false
	}

	item := q.Items[0]
	q.Items = q.Items[1:]
	return item, true
}

func (q *Queue) Front() (int, bool) {
	if q.IsEmpty() {
		return 0, false
	}
	return q.Items[0], true
}

func main() {
	queue := NewQueue()

	// Enqueue operations
	queue.Enqueue(10)
	queue.Enqueue(20)
	queue.Enqueue(30)
	queue.Display() // Queue: [10 20 30]

	// Dequeue operation
	item, _ := queue.Dequeue()
	fmt.Printf("Dequeued: %d\n", item) // Dequeued: 10
	queue.Display()                    // Queue: [20 30]

	// Front operation
	front, _ := queue.Front()
	fmt.Printf("Front element: %d\n", front) // Front element: 20
}
