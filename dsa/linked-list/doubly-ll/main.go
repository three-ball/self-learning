package main

import (
	"errors"
	"fmt"
)

var (
	ErrLinkedListEmpty        = errors.New("linked list is empty")
	ErrLinkedListNodeNotFound = errors.New("errors node not found")
)

// nodeData represents the data stored in each node
type nodeData struct {
	id   string // unique identifier for the node
	data int    // actual data value stored in the node
}

// newNodeData creates a new nodeData instance
func newNodeData(id string, data int) *nodeData {
	return &nodeData{
		id:   id,
		data: data,
	}
}

// Node represents a single node in the linked list
// Contains data and a pointer to the next node
type Node struct {
	data     *nodeData
	Next     *Node
	Previous *Node
}

// newNode creates a new Node instance
func newNode(n *nodeData, next *Node, previous *Node) *Node {
	return &Node{
		data:     n,
		Next:     next,
		Previous: previous,
	}
}

type DoublyLinkedList struct {
	head *Node
	tail *Node
	size int
}

func NewDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{
		head: nil,
		tail: nil,
		size: 0,
	}
}

// Size returns the current number of nodes in the linked list
// Time Complexity: O(1)
func (d *DoublyLinkedList) Size() int {
	return d.size
}

// IsEmpty checks if the linked list is empty
// Time Complexity: O(1)
func (d *DoublyLinkedList) IsEmpty() bool {
	return d.head == nil
}

// GetHead returns the head node (for external access if needed)
// Time Complexity: O(1)
func (d *DoublyLinkedList) GetHead() *Node {
	return d.head
}

// GetTail returns the tail node
func (d *DoublyLinkedList) GetTail() *Node {
	return d.tail
}

// Clear removes all nodes from the doubly linked list
// Time Complexity: O(1) - Go's garbage collector handles cleanup
func (d *DoublyLinkedList) Clear() {
	d.head = nil
	d.tail = nil
	d.size = 0
}

// InsertHead inserts a new node at the beginning of the doubly linked list
// The new node becomes the new head
// Time Complexity: O(1)
func (d *DoublyLinkedList) InsertHead(id string, data int) {
	nodeData := newNodeData(id, data)
	newNode := newNode(nodeData, d.head, nil)

	if !d.IsEmpty() {
		d.head.Previous = newNode
	} else {
		// first node is also the tail
		d.tail = newNode
	}

	d.head = newNode
	d.size++
}

// InsertTail inserts a new node at the end of the doubly linked list
// The new node becomes the new tail
// Time Complexity: O(1)
func (d *DoublyLinkedList) InsertTail(id string, data int) {
	nodeData := newNodeData(id, data)
	newNode := newNode(nodeData, nil, d.tail)

	if d.tail != nil {
		d.tail.Next = newNode
	} else {
		// First node is also the head
		d.head = newNode
	}

	d.tail = newNode
	d.size++
}

// Delete removes the first node with the specified id from the doubly linked list
// Returns an error if the list is empty or node is not found
// Time Complexity: O(n) - may need to traverse entire list
func (d *DoublyLinkedList) Delete(id string) error {
	if d.IsEmpty() {
		return ErrLinkedListEmpty
	}
	current := d.head
	for current != nil {
		if current.data.id == id {
			if current.Previous == nil {
				// head node
				d.head = current.Next
				d.head.Previous = nil
			} else {
				current.Previous.Next = current.Next
			}

			if current.Next == nil {
				// tail node
				d.tail = current.Previous
				d.tail.Next = nil
			} else {
				current.Next.Previous = current.Previous
			}

			d.size--
			return nil
		}
		current = current.Next
	}
	return ErrLinkedListNodeNotFound
}

func (d *DoublyLinkedList) Search(id string) (int, error) {
	if d.IsEmpty() {
		return 0, ErrLinkedListEmpty
	}

	current := d.head
	for current != nil {
		if current.data.id == id {
			return current.data.data, nil
		}

		current = current.Next
	}

	return 0, ErrLinkedListNodeNotFound
}

// Display prints the doubly linked list in forward direction
// Used for debugging and visualization
// Time Complexity: O(n)
func (d *DoublyLinkedList) Display() {
	if d.IsEmpty() {
		fmt.Println("List is empty")
		return
	}

	current := d.head
	fmt.Print("nil <- ")
	for current != nil {
		fmt.Printf("[%s:%d]", current.data.id, current.data.data)
		if current.Next != nil {
			fmt.Print(" <-> ")
		}
		current = current.Next
	}
	fmt.Println(" -> nil")
}

// DisplayReverse prints the doubly linked list in reverse direction
// Used for debugging and visualization of backward traversal
// Time Complexity: O(n)
func (d *DoublyLinkedList) DisplayReverse() {
	if d.IsEmpty() {
		fmt.Println("List is empty")
		return
	}

	current := d.tail
	fmt.Print("nil <- ")
	for current != nil {
		fmt.Printf("[%s:%d]", current.data.id, current.data.data)
		if current.Previous != nil {
			fmt.Print(" <-> ")
		}
		current = current.Previous
	}
	fmt.Println(" -> nil")
}

// Example usage and testing
func main() {
	dll := NewDoublyLinkedList()

	fmt.Println("=== Doubly Linked List Demo ===")

	// Insert elements
	dll.InsertHead("2", 20)
	dll.InsertHead("1", 10)
	dll.InsertTail("3", 30)
	dll.InsertTail("4", 40)

	fmt.Println("Forward display:")
	dll.Display()

	fmt.Println("Reverse display:")
	dll.DisplayReverse()

	fmt.Printf("Size: %d\n", dll.Size())

	// Search
	if val, err := dll.Search("3"); err == nil {
		fmt.Printf("Found node '3': %d\n", val)
	}

	// Delete
	dll.Delete("2")
	fmt.Println("After deleting '2':")
	dll.Display()
}
