package main

import "fmt"

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// IsEmpty checks if the stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

// Push adds an item onto the stack
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item from the stack
func (s *Stack[T]) Pop() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}

	lastIndex := len(s.items) - 1
	item := s.items[lastIndex]
	s.items = s.items[:lastIndex]
	return item, true
}

// Peek returns the top item without removing it
func (s *Stack[T]) Peek() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}

	lastIndex := len(s.items) - 1
	return s.items[lastIndex], true
}

func main() {
	// Test with strings
	fmt.Println("\n2. Testing String Stack:")
	stringStack := NewStack[string]()

	fmt.Println("Pushing: 'first', 'second', 'third'")
	stringStack.Push("first")
	stringStack.Push("second")
	stringStack.Push("third")

	fmt.Printf("Size: %d\n", stringStack.Size())

	if val, ok := stringStack.Peek(); ok {
		fmt.Printf("Top element: %s\n", val)
	}

	fmt.Println("Popping all strings:")
	for !stringStack.IsEmpty() {
		if val, ok := stringStack.Pop(); ok {
			fmt.Printf("Popped: '%s'\n", val)
		}
	}

	// Test LIFO behavior
	fmt.Println("\n3. Testing LIFO (Last In, First Out) Behavior:")
	lifoStack := NewStack[string]()

	operations := []string{"A", "B", "C", "D"}
	fmt.Printf("Pushing in order: %v\n", operations)
	for _, item := range operations {
		lifoStack.Push(item)
	}

	fmt.Print("Popping order: ")
	poppedItems := make([]string, 0)
	for !lifoStack.IsEmpty() {
		if val, ok := lifoStack.Pop(); ok {
			poppedItems = append(poppedItems, val)
		}
	}
	fmt.Printf("%v (should be reverse order)\n", poppedItems)

	// Test with custom struct
	fmt.Println("\n4. Testing Custom Struct Stack:")
	type Person struct {
		Name string
		Age  int
	}

	personStack := NewStack[Person]()

	personStack.Push(Person{Name: "Alice", Age: 30})
	personStack.Push(Person{Name: "Bob", Age: 25})
	personStack.Push(Person{Name: "Charlie", Age: 35})

	fmt.Printf("Stack size: %d\n", personStack.Size())

	fmt.Println("Popping persons:")
	for !personStack.IsEmpty() {
		if person, ok := personStack.Pop(); ok {
			fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
		}
	}
}
