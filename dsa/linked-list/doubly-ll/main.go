package main

import (
	"errors"
	"fmt"
)

var (
	ErrLinkedListEmpty        = errors.New("linked list is empty")
	ErrLinkedListNodeNotFound = errors.New("node not found")
)

type Node struct {
	Key   string // Assume key is string and unique, for example: uuid, email, etc.
	Value int
	Next  *Node
	Prev  *Node
}

type DLL struct {
	nodeDict map[string]*Node
	head     *Node
	tail     *Node
}

// NewDLL create the new doubly linked list
// with "sentinel" pattern: a dummy head and a dummy tail
func NewDLL() *DLL {
	head := &Node{}
	tail := &Node{}

	head.Next = tail
	tail.Prev = head

	return &DLL{
		nodeDict: make(map[string]*Node),
		head:     head,
		tail:     tail,
	}
}

// Size returns the number of nodes in the list
// Time Complexity: O(1)
func (dll *DLL) Size() int {
	return len(dll.nodeDict)
}

// IsEmpty checks if the list is empty
// Time Complexity: O(1)
func (dll *DLL) IsEmpty() bool {
	return dll.Size() == 0
}

// Append adds a new node at the end (before tail sentinel)
// Time Complexity: O(1)
func (dll *DLL) Append(key string, value int) {
	newNode := &Node{
		Key:   key,
		Value: value,
	}

	newNode.Prev = dll.tail.Prev
	newNode.Next = dll.tail

	dll.tail.Prev.Next = newNode
	dll.tail.Prev = newNode

	dll.nodeDict[key] = newNode
}

// Prepend adds a new node at the beginning (after head sentinel)
// Time Complexity: O(1)
func (dll *DLL) InsertHead(key string, value int) {
	newNode := &Node{
		Key:   key,
		Value: value,
	}

	newNode.Prev = dll.head
	newNode.Next = dll.head.Next

	dll.head.Next.Prev = newNode
	dll.head.Next = newNode

	dll.nodeDict[key] = newNode
}

// Remove deletes a node with the given key
// Time Complexity: O(1)
func (dll *DLL) Remove(key string) error {
	if dll.IsEmpty() {
		return ErrLinkedListEmpty
	}

	if node, ok := dll.nodeDict[key]; ok {
		node.Next.Prev = node.Prev
		node.Prev.Next = node.Next
		return nil
	}

	return ErrLinkedListNodeNotFound
}

// Search finds a node with the given key and returns its value
// Time Complexity: O(1)
func (dll *DLL) Search(key string) (int, error) {
	if dll.IsEmpty() {
		return 0, ErrLinkedListEmpty
	}

	if node, ok := dll.nodeDict[key]; ok {
		return node.Value, nil
	}

	return 0, ErrLinkedListNodeNotFound
}

// Update modifies the value of an existing key
// Time Complexity: O(1)
func (dll *DLL) Update(key string, value int) error {
	if node, ok := dll.nodeDict[key]; ok {
		node.Value = value
		return nil
	}
	return ErrLinkedListNodeNotFound
}

// Display prints the list from head to tail
// Time Complexity: O(n)
func (dll *DLL) Display() {
	if dll.IsEmpty() {
		fmt.Println("List is empty")
		return
	}

	current := dll.head.Next
	fmt.Print("HEAD <-> ")
	for current != dll.tail {
		fmt.Printf("[%s:%d] <-> ", current.Key, current.Value)
		current = current.Next
	}
	fmt.Println("TAIL")
}

// Clear removes all nodes from the list
// Time Complexity: O(1)
func (dll *DLL) Clear() {
	dll.head.Next = dll.tail
	dll.tail.Prev = dll.head
	dll.nodeDict = make(map[string]*Node)
}

// Example usage
func main() {
	dll := NewDLL()

	fmt.Println("=== Doubly Linked List with HashMap Demo ===")

	// Test operations
	dll.Append("first", 10)
	dll.Append("second", 20)
	dll.InsertHead("zero", 0)
	dll.Append("third", 30)

	fmt.Println("Forward display:")
	dll.Display()

	// Search operations
	if val, err := dll.Search("second"); err == nil {
		fmt.Printf("Found 'second': %d\n", val)
	}

	// Update
	dll.Update("second", 25)
	fmt.Println("After updating 'second' to 25:")
	dll.Display()

	// Remove
	dll.Remove("zero")
	fmt.Println("After removing 'zero':")
	dll.Display()

	fmt.Printf("Size: %d\n", dll.Size())
}
