package ast

// Stack of character allow to push and pop character from a string
type CharStack struct {
	data string
}

// Push a character on the stack
func (s *CharStack) Push(c byte) {
	s.data = string(c) + s.data
}

// Push string push a string on the stack
func (s *CharStack) PushString(str string) {
	s.data = str + s.data
}

// Pop a character from the stack
func (s *CharStack) Pop() byte {
	if len(s.data) == 0 {
		return 0
	}
	c := s.data[0]
	s.data = s.data[1:]
	return c
}

// PopVerify pops a character from the stack and verify it
func (s *CharStack) PopVerify(c byte) bool {
	if len(s.data) == 0 {
		return false
	}
	if s.data[0] != c {
		return false
	}
	s.data = s.data[1:]
	return true
}

// PopVerifyString pops a string from the stack and verify it
func (s *CharStack) PopVerifyString(str string) bool {
	if len(s.data) < len(str) {
		return false
	}
	if s.data[:len(str)] != str {
		return false
	}
	s.data = s.data[len(str):]
	return true
}

// PopString pops a string from the stack
func (s *CharStack) PopString(n int) string {
	if len(s.data) < n {
		return ""
	}
	str := s.data[:n]
	s.data = s.data[n:]
	return str
}

// Peek at the top of the stack
func (s *CharStack) Peek() byte {
	if len(s.data) == 0 {
		return 0
	}
	return s.data[0]
}

// IsEmpty returns true if the stack is empty
func (s *CharStack) IsEmpty() bool {
	return len(s.data) == 0
}

// String returns the string representation of the stack
func (s *CharStack) String() string {
	return s.data
}

// Len returns the length of the stack
func (s *CharStack) Len() int {
	return len(s.data)
}
