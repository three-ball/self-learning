package main

import (
	"fmt"
)

type KeyValue struct {
	Key   string
	Value interface{}
	Next  *KeyValue // for chainning collision
}

type HashTable struct {
	buckets    []*KeyValue
	size       int // size of the hash table
	count      int // number of elements in the hash table
	loadFactor float64
}

func NewHashTable(size int, loadFactor float64) *HashTable {
	return &HashTable{
		buckets:    make([]*KeyValue, size),
		size:       size,
		count:      0,
		loadFactor: loadFactor,
	}
}

// Hash is the hash function of HashTable
func (ht *HashTable) Hash(key string) int {
	hash := 5381
	for _, char := range key {
		hash = ((hash << 5) + hash) + int(char)
	}
	return hash % ht.size
}

func (ht *HashTable) PutWithoutChaining(key string, value interface{}) {
	// Without chaining, we simply insert the key-value pair
	hash := ht.Hash(key)
	ht.buckets[hash] = &KeyValue{
		Key:   key,
		Value: value,
		Next:  nil,
	}
	ht.count++
}

func (ht *HashTable) Put(key string, value interface{}) {
	hash := ht.Hash(key)
	if ht.buckets[hash] == nil {
		ht.buckets[hash] = &KeyValue{
			Key:   key,
			Value: value,
			Next:  nil,
		}
		ht.count++
		return
	}

	// Handle collision (chaining)
	current := ht.buckets[hash]
	for current != nil {
		if current.Key == key {
			current.Value = value
			return
		}

		if current.Next == nil {
			break
		}

		current = current.Next
	}

	current.Next = &KeyValue{
		Key:   key,
		Value: value,
		Next:  nil,
	}
	ht.count++

	if float64(ht.count)/float64(ht.size) > ht.loadFactor {
		ht.resize()
	}
}

func (ht *HashTable) Get(key string) (interface{}, bool) {
	hashKey := ht.Hash(key)
	// search the bucket by hash key
	current := ht.buckets[hashKey]

	// search the real data in the linked list
	for current != nil {
		if current.Key == key {
			return current.Value, true
		}
		current = current.Next
	}
	return nil, false
}

func (ht *HashTable) Delete(key string) bool {
	hashKey := ht.Hash(key)
	current := ht.buckets[hashKey]
	if current == nil {
		return false
	}

	previous := current
	for current != nil {
		if current.Key == key {
			previous.Next = current.Next
			ht.count--
			return true
		}
		previous = current
		current = current.Next
	}
	return false
}

func (ht *HashTable) Keys() []string {
	keys := make([]string, 0, ht.count)
	for _, bucket := range ht.buckets {
		current := bucket
		for current != nil {
			keys = append(keys, current.Key)
			current = current.Next
		}
	}
	return keys
}

func (ht *HashTable) Display() {
	fmt.Printf("Hash Table (Size: %d, Count: %d, Load Factor Treshold: %.2f, Current Load Factor: %.2f)\n",
		ht.size,
		ht.count,
		ht.loadFactor,
		float64(ht.count)/float64(ht.size))

	for i, bucket := range ht.buckets {
		fmt.Printf("Bucket %d: \n", i)
		current := bucket
		if current == nil {
			fmt.Println("-> Empty")
		} else {
			for current != nil {
				fmt.Printf("-> [Key: %s, Value: %v]\n", current.Key, current.Value)
				current = current.Next
			}
		}
		fmt.Println("-------")
	}
}

func (ht *HashTable) resize() {
	oldBuckets := ht.buckets
	oldSize := ht.size

	ht.size = oldSize * 2
	ht.buckets = make([]*KeyValue, ht.size)

	// Rehash all existing keys
	for _, bucket := range oldBuckets {
		current := bucket
		for current != nil {
			ht.Put(current.Key, current.Value)
			current = current.Next
		}
	}
}

func main() {
	ht := NewHashTable(10, 0.8)

	// Insert key-value pairs
	ht.Put("name", "John")
	ht.Put("age", 30)
	ht.Put("city", "New York")
	ht.Put("country", "USA")
	ht.Put("occupation", "Engineer")
	ht.Put("hobby", "Photography")
	ht.Put("language", "Go")
	ht.Put("framework", "Gin")

	ht.Display()

	// Get values
	if value, exists := ht.Get("name"); exists {
		fmt.Printf("name: %v\n", value)
	}

	// Update value
	ht.Put("age", 31)
	age, res := ht.Get("age")
	if res {
		fmt.Printf("Updated age: %v\n", age)
	}

	ht.Display()

	// Delete key
	ht.Delete("city")

	// Show all keys
	fmt.Printf("All keys: %v\n", ht.Keys())
}
