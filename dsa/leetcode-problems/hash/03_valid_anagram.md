# Valid Anagram

## Problem Statement

Given two strings `s` and `t`, return `true` if `t` is an anagram of `s`, and `false` otherwise.

## First Attempt

- First attempt is failed. The code only checks if all characters in `t` are in `s`, but does not check the counts of each character.

```go
func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	var sDict = make(map[rune]bool)

	for _, value := range s {
		sDict[value] = true
	}

	for _, value := range t {
		if _, ok := sDict[value]; !ok {
			return false
		}
	}
	return true
}
```

## Second Attempt

- Use two hash maps to count the occurrences of each character in both strings, then compare the two maps. Simply is 3 loops.
- Time complexity is `O(n)`, space complexity is `O(1)` since the size of the hash map is bounded by the number of unique characters (which is constant for a fixed character set).

```go
func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	var sDict = make(map[rune]int)
	var tDict = make(map[rune]int)
	for _, value := range s {
		sDict[value] += 1
	}

	for _, value := range t {
		tDict[value] += 1
	}

	for k, v := range sDict {
		if tDict[k] != v {
			return false
		}
	}

	return true
}
```

![valid_anagram_01](images/valid_anagram_01.png)

## Third Attempt

- We can do better by using only one hash map. First, count the occurrences of each character in `s`, then decrement the count for each character in `t`. If any count goes negative or if there are any non-zero counts left at the end, return false.
- Time complexity is `O(n)`, space complexity is `O(1)`.

```go
	var lenS = len(s)
	if lenS != len(t) {
		return false
	}

	var sDict = make(map[byte]int)
	for index := 0; index < lenS; index++ {
		sDict[s[index]]++
		sDict[t[index]]--
	}

	for _, v := range sDict {
		if v != 0 {
			return false
		}
	}

	return true
```

![valid_anagram_02](images/valid_anagram_02.png)

- Still not better.

## Fourth Attempt

- We can do better by using an array instead of a hash map. Since the problem states that the strings consist of lowercase English letters, we can use an array of size 26 to count the occurrences of each character.
- Time complexity is `O(n)`, space complexity is `O(1)`.

```go
	var lenS = len(s)
	if lenS != len(t) {
		return false
	}

	var sDict = make([]int, 26)
	for index := 0; index < lenS; index++ {
		sIndex := s[index] - 'a'
		tIndex := t[index] - 'a'
		sDict[sIndex]++
		sDict[tIndex]--
	}

	for _, v := range sDict {
		if v != 0 {
			return false
		}
	}

	return true

// map version, still solve the problem but slower performance than array version
func isAnagram3(s string, t string) bool {
	var lenS = len(s)
	if lenS != len(t) {
		return false
	}

	var sDict = make(map[byte]int)
	for index := 0; index < lenS; index++ {
		sDict[s[index]]++
		sDict[t[index]]--
	}

	for _, v := range sDict {
		if v != 0 {
			return false
		}
	}

	return true
}
```

![valid_anagram_03](images/valid_anagram_03.png)

- Key point: Some special cases can be optimized by using arrays instead of hash maps. I think if we use map solution, it will fit more general cases, but if we know the input is limited to certain characters, we can use arrays for better performance.