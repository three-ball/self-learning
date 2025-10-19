# 20. Valid Parentheses

## Statement

Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.

An input string is valid if:

1. Open brackets must be closed by the same type of brackets.
2. Open brackets must be closed in the correct order.
3. Every close bracket has a corresponding open bracket of the same type.
 
## Solution in Go

- The opening brackets are pushed onto a stack.
- When a closing bracket is encountered, the top element is popped from the stack and checked for a match.
- If there is a mismatch or if the stack is empty when a closing bracket is encountered, the string is invalid.
- Finally, if the stack is empty after processing the entire string, the string is valid.

```go
const OPENBRACKET = -1
const CLOSEBRACKET = 1

var bracketCharDict = map[rune]int{
	'(': OPENBRACKET,
	'{': OPENBRACKET,
	'[': OPENBRACKET,
	')': CLOSEBRACKET,
	'}': CLOSEBRACKET,
	']': CLOSEBRACKET,
}

var bracketCharPairs = map[rune]rune{
	'(': ')',
	'{': '}',
	'[': ']',
}

type stack struct {
	items []rune
}

func newStack() stack {
	return stack{
		items: make([]rune, 0),
	}
}

func (s *stack) pop() (rune, bool) {
	if len(s.items) == 0 {
		return 0, false
	}

	lastIndex := len(s.items) - 1
	item := s.items[lastIndex]
	s.items = s.items[:lastIndex]
	return item, true
}

func (s *stack) push(a rune) {
	s.items = append(s.items, a)
}

func initParentheseStacks() stack {
	return newStack()
}

func isValid(s string) bool {
	openStack := initParentheseStacks()

	for _, char := range s {
		if val, ok := bracketCharDict[char]; ok {
			if !ok {
				continue
			}

			if val == OPENBRACKET {
				openStack.push(char)
			} else {
				item, ok := openStack.pop()
				if !ok {
					return false
				}

				if bracketCharPairs[item] != char {
					return false
				}
			}
		}
	}

	return len(openStack.items) == 0
}
```