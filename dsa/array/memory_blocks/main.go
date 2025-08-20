package main

import (
	"fmt"
)

// Given an array representing memory blocks [0,1,0,0,1,0,1,0] where:
// - 0 = free block
// - 1 = allocated block

// Tasks:
// 1. Find the largest contiguous free memory segment
// 2. Allocate N consecutive blocks (return starting index or -1 if impossible)
// 3. Deallocate blocks from index i to j

func findLargestFreeSegment(memory []int) int {
	var currentLength = 0
	var maxLength = 0
	for _, value := range memory {
		if value == 0 {
			currentLength++
			if currentLength > maxLength {
				maxLength = currentLength
			}
		} else {
			currentLength = 0
		}
	}
	return maxLength
}

func allocateNConsecutiveBlockIndex(memory []int, n int) int {
	if n <= 0 {
		return -1
	}

	totalFree := 0
	for index, _ := range memory {
		totalFree++
		if totalFree == n {
			// index  = [0, 1, 2, 3, 4, 5]
			// memory = [0, 1, 1, 0, 0, 0]
			// n = 2
			// index = 4 match the condition, return index = 3 -> 4 - n + 1
			return index - n + 1
		} else {
			totalFree = 0
		}
	}

	return totalFree
}

func deallocateRange(memory []int, from int, to int) []int {
	maxIndex := len(memory) - 1
	if from > maxIndex || to > maxIndex {
		return memory // do nothing
	}

	if from < 0 || to < 0 {
		return memory // do nothing
	}

	if from > to {
		return memory // do nothing
	}
	for index := from; index <= to; index++ {
		memory[index] = 0
	}

	return memory
}

func main() {
	memory := []int{0, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1}
	largestSegment := findLargestFreeSegment(memory)
	println("Largest contiguous free memory segment:", largestSegment)
	allocateIndex := allocateNConsecutiveBlockIndex(memory, 6)
	println("allocateIndex:", allocateIndex)
	deallocMemory := deallocateRange(memory, 3, 4)
	fmt.Println("deallocateIndex", deallocMemory)
}
