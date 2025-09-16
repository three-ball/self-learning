# 146. LRU Cache

## Statement

> Design a data structure that follows the constraints of a Least Recently Used (LRU) cache.

- Least Recently Used (LRU) cache: Discards least recently used items first.

```go
package main

type Node struct {
	Key      int // Store the key, not value
	Value    int // Store the value
	Next     *Node
	Previous *Node
}

type LRUCache struct {
	capacity int           // capacity of the cache, as the statement requires
	cache    map[int]*Node // statement: "The functions get and put must each run in O(1) average time complexity." and "key-value pairs"
	// tail is the most recently used node
	tail *Node
	// head is the least recently used node
	head *Node
}

func Constructor(capacity int) LRUCache {
	// dummy node
	head := &Node{}
	tail := &Node{}

	head.Next = tail
	tail.Previous = head
	return LRUCache{
		capacity: capacity,
		cache:    make(map[int]*Node, capacity),
		tail:     tail,
		head:     head,
	}
}

func (this *LRUCache) insertToTail(n *Node) {
	n.Previous = this.tail.Previous
	n.Next = this.tail
	this.tail.Previous.Next = n
	this.tail.Previous = n
	// this.tail = n
}

func (this *LRUCache) deleteNode(node *Node) {
	node.Previous.Next = node.Next
	node.Next.Previous = node.Previous
}

func (this *LRUCache) moveToTail(n *Node) {
	// this.deleteNode(n)
	// this.insertToTail(n)
	n.Previous.Next = n.Next
	n.Next.Previous = n.Previous

	n.Previous = this.tail.Previous
	n.Next = this.tail
	this.tail.Previous.Next = n
	this.tail.Previous = n
}

func (this *LRUCache) Get(key int) int {
	if node, exists := this.cache[key]; exists {
		this.moveToTail(node)
		return node.Value
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if node := this.cache[key]; node != nil {
		// if exist, update and move to tail
		node.Key = key
		node.Value = value
		this.moveToTail(node)
	} else {
		node := &Node{
			Key:   key,
			Value: value,
		}

		if len(this.cache) == this.capacity {
			oldestNode := this.head.Next
			this.deleteNode(oldestNode)
			delete(this.cache, oldestNode.Key)

		}

		this.insertToTail(node)
		this.cache[key] = node
	}

	return
}
```