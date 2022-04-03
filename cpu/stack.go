package cpu

// Stack is a LIFO of uint16 values
type Stack struct {
	entries []uint16
	size   uint16
}

func NewStack(size uint16) *Stack {
	return &Stack{
		entries: make([]uint16, 0),
		size: size,
	}
}

func (s *Stack) Size() int {
	return len(s.entries)
}

func (s *Stack) MaxSize() int {
	return int(s.size)
}

// Push pushes a value onto the stack if possible
func (s *Stack) Push(value uint16) bool {
	if len(s.entries) >= int(s.size) {
		return false
	}
	s.entries = append(s.entries, value)
	return true
}

// Pop pops a value off the stack
func (s *Stack) Pop() (uint16, bool) {
	if len(s.entries) == 0 {
		return 0, false
	}
	last := len(s.entries) - 1
	value := s.entries[last]
	s.entries = s.entries[:last]
	return value, true
}