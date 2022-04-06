package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Opcode_returns_name_using_base_struct(t *testing.T) {
	op := OxClearScreen{Opcode{0x0, "Test Clear"}}
	assert.Equal(t, "Test Clear", op.Name())
}
