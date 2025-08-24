package main

import "fmt"

type Node struct {
	data int
	next *Node
}

type LinkedQueue struct {
	front *Node
	rear  *Node
	size  int
}

func NewLinkedQueue() *LinkedQueue {
	return &LinkedQueue{
		front: nil,
		rear:  nil,
		size:  0,
	}
}

func (q *LinkedQueue) Display() {
	current := q.front
	for current != nil {
		fmt.Print(current.data, " ")
		current = current.next
	}
	fmt.Println()
}

func (q *LinkedQueue) IsEmpty() bool {
	return q.size == 0
}

func (q *LinkedQueue) Size() int {
	return q.size
}

func (q *LinkedQueue) Enqueue(item int) {
	newNode := &Node{data: item, next: nil}

	// empty queue
	if q.IsEmpty() {
		// front → [1] ← rear
		// 		    ↑
		//      Same node!
		q.front = newNode
		q.rear = newNode
	} else {
		// assign the "next" pointer of the current rear to the new node
		// 	front → [1] → [2]
		//     ↑     ↑
		//   rear  newNode
		q.rear.next = newNode

		// update the rear pointer to the new node
		// front → [1] → [2] ← rear
		q.rear = newNode
	}
	q.size++
}

func (q *LinkedQueue) Dequeue() (int, bool) {
	if q.IsEmpty() {
		return 0, false
	}

	it := q.front.data
	q.front = q.front.next
	q.size--

	if q.IsEmpty() {
		// NOW the queue is empty, but rear still points to the old node!
		// rear → [Node{data:10, next:nil}] (garbage/dangling pointer)
		// front → nil
		q.rear = nil
	}

	return it, true
}

func main() {
	queue := NewLinkedQueue()

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	queue.Display()

	for !queue.IsEmpty() {
		if item, ok := queue.Dequeue(); ok {
			fmt.Println(item)
		}
		queue.Display()
	}
}
