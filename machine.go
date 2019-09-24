package machine

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
)

var (
	MemoryOverflowError = errors.New("memory overflow")
	CycleError          = errors.New("'[' or ']' expected, got nothing")
)

type Machine struct {
	mem     [30000]byte
	pointer int
}

func Execute(reader io.Reader) ([]byte, error) {
	instructions, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var (
		m      Machine
		result []byte
	)

	if bytes.Count(instructions, []byte{'['}) != bytes.Count(instructions, []byte{']'}) {
		return nil, CycleError
	}

	for i := 0; i < len(instructions); i++ {
		switch instructions[i] {
		case '>':
			if m.pointer += 1; m.pointer > len(m.mem)-1 {
				return nil, MemoryOverflowError
			}
		case '<':
			if m.pointer -= 1; m.pointer < 0 {
				return nil, MemoryOverflowError
			}
		case '+':
			m.mem[m.pointer]++
		case '-':
			m.mem[m.pointer]--
		case '.':
			result = append(result, m.mem[m.pointer])
		case ',':
			// TODO
		case '[':
			if m.mem[m.pointer] != 0 {
				break
			}

			for depth := 1; depth > 0; {
				if i += 1; i > len(instructions)-1 {
					return nil, CycleError
				}

				switch instructions[i] {
				case '[':
					depth++
				case ']':
					depth--
				}
			}
		case ']':
			for depth := 1; depth > 0; {
				if i -= 1; i < 0 {
					return nil, CycleError
				}

				switch instructions[i] {
				case '[':
					depth--
				case ']':
					depth++
				}
			}
			i--
		}
	}

	return result, nil
}
