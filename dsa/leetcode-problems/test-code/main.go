package main

// wrong answer, what if more than 1 char
// s = "aacc"
// t= "ccac"
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

func isAnagram2(s string, t string) bool {
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

func isAnagram4(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	var sDict = make(map[rune]int)
	for _, value := range s {
		sDict[value]++
	}

	for _, value := range t {
		if sDict[value] == 0 {
			return false
		}
		sDict[value]--
	}

	return true
}

func isAnagram5(s string, t string) bool {
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
}
