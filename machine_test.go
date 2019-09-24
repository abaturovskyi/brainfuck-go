package machine

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

const (
	basicTest = `
		+++++++++++++++++++++++++++++++++++++++++++++
		+++++++++++++++++++++++++++.+++++++++++++++++
		++++++++++++.+++++++..+++.-------------------
		---------------------------------------------
		---------------.+++++++++++++++++++++++++++++
		++++++++++++++++++++++++++.++++++++++++++++++
		++++++.+++.------.--------.------------------
		---------------------------------------------
		----.`

	complexTest = `++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.`

	overflowTest = `>+++<<`

	unclosedCycleTest1 = `>++[[-]`
	unclosedCycleTest2 = `>++[-]]`
	unclosedCycleTest3 = `>++][-]`
	unclosedCycleTest4 = `>++[-][`
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name   string
		reader io.Reader
		exp    string
		err    error
	}{
		{
			name:   "basic",
			reader: bytes.NewBufferString(basicTest),
			exp:    "Hello World!",
		},
		{
			name:   "complex",
			reader: bytes.NewBufferString(complexTest),
			exp:    "Hello World!",
		},
		{
			name:   "memory_overflow",
			reader: bytes.NewBufferString(overflowTest),
			err:    MemoryOverflowError,
		},
		{
			name:   "cycle_1",
			reader: bytes.NewBufferString(unclosedCycleTest1),
			err:    CycleError,
		},
		{
			name:   "cycle_2",
			reader: bytes.NewBufferString(unclosedCycleTest2),
			err:    CycleError,
		},
		{
			name:   "cycle_3",
			reader: bytes.NewBufferString(unclosedCycleTest3),
			err:    CycleError,
		},
		{
			name:   "cycle_4",
			reader: bytes.NewBufferString(unclosedCycleTest4),
			err:    CycleError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Execute(tt.reader)
			if tt.err != err {
				t.Fatal(err)
			}

			if !bytes.Equal([]byte(tt.exp), result) {
				t.Error(fmt.Sprintf("expected: '%s', got: '%s'", tt.exp, string(result)))
			}
		})
	}
}
