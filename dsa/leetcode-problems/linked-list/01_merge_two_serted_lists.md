# Merge Two Sorted Lists

## Problem Statement

You are given the heads of two sorted linked lists list1 and list2.
Merge the two lists into one sorted list. The list should be made by splicing together the nodes of the first two lists.
Return the head of the merged linked list.

## Solution

- Create a dummy node to serve as the starting point of the merged list.
- This dummy node is "hold" the head of the merged list. Then we use the `tem[` variable to start building the merged list.
- We loop through both lists, comparing the current nodes of each list. The smaller node is appended to the merged list, and we move the pointer of that list forward until the current node of this list is bigger than the other list's current node.
- The shorter lists will be exhausted first, we then append the remaining nodes of the longer list to the merged list.

```go
type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	var dummyList = &ListNode{} // dummy hold the head

	temp := dummyList

	p1, p2 := list1, list2

	for p1 != nil && p2 != nil {
		if p1.Val > p2.Val {
			temp.Next = p2
			p2 = p2.Next
		} else {
			temp.Next = p1
			p1 =  p1.Next
		}

		temp = temp.Next
	}

	if p1 != nil {
		temp.Next = p1
	} else {
		temp.Next = p2
	}

	return dummyList.Next
}
```