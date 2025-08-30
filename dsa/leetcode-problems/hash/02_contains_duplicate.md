# Contains Duplicate

## Problem Statement

Given an integer array nums, return true if any value appears at least twice in the array, and return false if every element is distinct.

## Submits

### First Attempt

- First attempt may be brute force, check every pair of elements. But we can do better. Using a hash map, we can track the elements we have seen so far.

```go
func containsDuplicate(nums []int) bool {
	var distinctDict = make(map[int]bool)
	for _, value := range nums {
		if distinctDict[value] {
			return true
		}
		distinctDict[value] = true
	}
	return false
}
```
![contain_duplicate_hash_01](images/contain_duplicate_hash_01.png)

- Why using hash but it's still beat only 55% of Go submissions?

### Second Attempt

- We can do better by checking if the key exists in the map instead of checking its value.

```go
func containsDuplicate(nums []int) bool {
	var distinctDict = make(map[int]bool)
	for _, value := range nums {
		if _, ok := distinctDict[value];ok {
			return true
		}
		distinctDict[value] = true
	}
	return false
}
```

![alt text](image.png/contain_duplicate_hash_02.png)

- Best practice: use `if _, ok := map[key]; ok {}` to check if a key exists in a map.