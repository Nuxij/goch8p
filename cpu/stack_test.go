package cpu

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStack(t *testing.T) {
	type args struct {
		size uint16
	}
	tests := []struct {
		name string
		args args
		want *Stack
	}{
		{"large", args{0x1000}, &Stack{entries: make([]uint16, 0), size: 0x1000}},
		{"small", args{0x100}, &Stack{entries: make([]uint16, 0), size: 0x100}},
		{"chip8", args{0x10}, &Stack{entries: make([]uint16, 0), size: 0x10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStack(tt.args.size)
			assert.EqualValues(t, tt.args.size, got.MaxSize())
			assert.EqualValues(t, tt.want.size, got.MaxSize())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStack() = %v, want %v", got, tt.want)
			}	
		})
	}
}

func TestStack_Push(t *testing.T) {
	type args struct {
		value uint16
	}
	tests := []struct {
		name     string
		s        *Stack
		args     args
		want     bool
		wantSize int
	}{
		{"empty", &Stack{entries: make([]uint16, 0), size: 0x100}, args{0x100}, true, 1},
		{"full", &Stack{entries: make([]uint16, 0x100), size: 0x100}, args{0x100}, false, 256},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Push(tt.args.value); got != tt.want {
				t.Errorf("Stack.Push() = %v, want %v", got, tt.want)
			}
			if tt.s.Size() != tt.wantSize {
				t.Errorf("Stack.Size() = %v, want %v", tt.s.Size(), tt.wantSize)
			}
		})
	}
}

func TestStack_Pop(t *testing.T) {
	tests := []struct {
		name  string
		s     *Stack
		want  uint16
		want1 bool
	}{
		{"empty", &Stack{entries: make([]uint16, 0), size: 0x100}, 0, false},
		{"full", &Stack{entries: []uint16{0x100, 0x200}, size: 0x100}, 0x200, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.Pop()
			if got != tt.want {
				t.Errorf("Stack.Pop() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Stack.Pop() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
