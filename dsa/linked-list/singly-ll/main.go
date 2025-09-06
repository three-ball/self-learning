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
	data *nodeData
	Next *Node
}

// newNode creates a new Node instance
func newNode(n *nodeData, next *Node) *Node {
	return &Node{
		data: n,
		Next: next,
	}
}

// SinglyLinkedList represents a singly linked list data structure
// Only maintains a reference to the head node and tracks size
type SinglyLinkedList struct {
	head *Node
	size int
}

// NewSinglyLinkedList creates and returns a new empty singly linked list
func NewSinglyLinkedList() *SinglyLinkedList {
	return &SinglyLinkedList{
		head: nil,
		size: 0,
	}
}

// Size returns the current number of nodes in the linked list
// Time Complexity: O(1)
func (s *SinglyLinkedList) Size() int {
	return s.size
}

// IsEmpty checks if the linked list is empty
// Time Complexity: O(1)
func (s *SinglyLinkedList) IsEmpty() bool {
	return s.head == nil
}

// GetHead returns the head node (for external access if needed)
// Time Complexity: O(1)
func (s *SinglyLinkedList) GetHead() *Node {
	return s.head
}

// Clear removes all nodes from the linked list
// Time Complexity: O(1) - Go's garbage collector handles cleanup
func (s *SinglyLinkedList) Clear() {
	s.head = nil
	s.size = 0
}

// Display prints the linked list in a readable format
// Used for debugging and visualization
// Time Complexity: O(n)
func (s *SinglyLinkedList) Display() {
	if s.head == nil {
		fmt.Println("List is empty")
		return
	}

	current := s.head
	fmt.Print("Head -> ")
	for current != nil {
		fmt.Printf("[%s:%d] -> ", current.data.id, current.data.data)
		current = current.Next
	}
	fmt.Println("nil")
}

// InsertHead inserts a new node at the beginning of the linked list
// The new node becomes the new head
// Time Complexity: O(1)
func (s *SinglyLinkedList) InsertHead(id string, data int) {
	nodeData := newNodeData(id, data)
	node := newNode(nodeData, s.head)
	s.head = node
	s.size++
}

// InsertTail inserts a new node at the end of the linked list
// If list is empty, the new node becomes the head
// Time Complexity: O(n) - need to traverse to the end
func (s *SinglyLinkedList) InsertTail(id string, data int) {
	nodeData := newNodeData(id, data)

	// first node check (linked list empty)
	if s.IsEmpty() {
		s.head = newNode(nodeData, nil)
		s.size++
		return
	}

	currentNode := s.head
	for currentNode.Next != nil {
		currentNode = currentNode.Next
	}

	// code reach here, mean currentNode is now called "tail node"
	currentNode.Next = newNode(nodeData, nil)
	s.size++
}

// Delete removes the first node with the specified id from the linked list
// Returns an error if the list is empty or node is not found
// Time Complexity: O(n) - may need to traverse entire list
func (s *SinglyLinkedList) Delete(id string) error {
	if s.IsEmpty() {
		return ErrLinkedListEmpty
	}

	if s.head.data.id == id {
		s.head = s.head.Next
		s.size--
		return nil
	}

	currentNode := s.head
	for currentNode.Next != nil {
		if currentNode.Next.data.id == id {
			currentNode.Next = currentNode.Next.Next
			s.size--
			return nil
		}
		currentNode = currentNode.Next
	}
	return ErrLinkedListNodeNotFound
}

// Search finds a node with the specified id and returns its data value
// Returns an error if the list is empty or node is not found
// Time Complexity: O(n) - may need to traverse entire list
func (s *SinglyLinkedList) Search(id string) (int, error) {
	if s.IsEmpty() {
		return 0, ErrLinkedListEmpty
	}

	currentNode := s.head
	for currentNode != nil {
		if currentNode.data.id == id {
			return currentNode.data.data, nil
		}
		currentNode = currentNode.Next
	}
	return 0, ErrLinkedListNodeNotFound
}

func main() {
	// Create a new singly linked list
	ll := NewSinglyLinkedList()

	fmt.Println("=== Singly Linked List Demo ===")

	// Insert elements
	ll.InsertHead("1", 10)
	ll.InsertTail("2", 20)
	ll.InsertTail("3", 30)
	ll.InsertHead("0", 5)

	fmt.Println("After insertions:")
	ll.Display()

	// Search for an element
	if val, err := ll.Search("2"); err == nil {
		fmt.Printf("Found node with id '2': %d\n", val)
	} else {
		fmt.Printf("Search error: %v\n", err)
	}

	// Display list information
	fmt.Printf("List size: %d\n", ll.Size())
	fmt.Printf("Is empty: %v\n", ll.IsEmpty())

	// Delete an element
	if err := ll.Delete("2"); err == nil {
		fmt.Println("Successfully deleted node with id '2'")
	} else {
		fmt.Printf("Delete error: %v\n", err)
	}

	fmt.Println("After deletion:")
	ll.Display()
	fmt.Printf("New size: %d\n", ll.Size())
}
