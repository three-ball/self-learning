package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1 == nil {
		if list2 != nil {
			return list2
		}
		return nil
	}

	if list2 == nil {
		return list1
	}

	var dummyList = &ListNode{}

	currentL1 := list1
	for currentL1 != nil {
		currentL2 := list2
		for currentL2 != nil {
		}
	}

	return current
}
