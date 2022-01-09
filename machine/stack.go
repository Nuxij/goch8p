package machine

// Stack is a 16-element uint16 array
// uses the last element as the top of the stack
type Stack [17]uint16

// Push pushes a value onto the stack
func (s *Stack) Push(v uint16) bool {
	if s[16] == 15 {
		return false
	} else {
		s[16]++
		s[s[16]] = v
		return true
	}
}

// Pop pops a value off the stack
func (s *Stack) Pop() (uint16, bool) {
	if s[16] == 0 {
		return 0, false
	} else {
		s[16]--
		return s[s[16]], true
	}
}