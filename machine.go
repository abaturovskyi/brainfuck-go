package machine

import (
	"errors"
	"io"
	"io/ioutil"
)

var (
	MemoryOverflowError = errors.New("memory overflow")
	CycleError          = errors.New("'[' or ']' expected, got nothing")

	handlersSet = map[byte]func(m *Machine, s *scope) error{
		'>': func(m *Machine, s *scope) error {
			if m.pointer += 1; m.pointer > len(m.memory)-1 {
				return MemoryOverflowError
			}
			return nil
		},
		'<': func(m *Machine, s *scope) error {
			if m.pointer -= 1; m.pointer < 0 {
				return MemoryOverflowError
			}
			return nil
		},
		'+': func(m *Machine, s *scope) error {
			m.memory[m.pointer]++
			return nil
		},
		'-': func(m *Machine, s *scope) error {
			m.memory[m.pointer]--
			return nil
		},
		'.': func(m *Machine, s *scope) error {
			_, err := m.output.Write([]byte{m.memory[m.pointer]})
			return err
		},
		',': func(m *Machine, s *scope) error {
			// TODO
			return nil
		},
		'[': func(m *Machine, s *scope) error {
			// Seek to next ']' rune counting depth
			seekCursor := s.cursor
			for depth := 1; depth > 0; {
				if seekCursor += 1; seekCursor > len(s.instructions)-1 {
					return CycleError
				}

				switch s.instructions[seekCursor] {
				case '[':
					depth++
				case ']':
					depth--
				}
			}

			// Skip cycle condition
			if m.memory[m.pointer] == 0 {
				s.cursor = seekCursor
				return nil
			}

			// Fall into cycle, create new scope
			if err := m.run(&scope{instructions: s.instructions[s.cursor+1 : seekCursor]}); err != nil {
				return err
			}

			// Reset cursor, repeat cycle
			s.cursor--
			return nil
		},
		']': func(m *Machine, s *scope) error {
			return CycleError
		},
	}
)

type scope struct {
	instructions []byte
	cursor       int
}

type Machine struct {
	// 1 KB of memory
	memory  [1 << 10]byte
	pointer int

	output   io.Writer
	handlers map[byte]func(m *Machine, s *scope) error
}

func (m *Machine) run(s *scope) error {
	for s.cursor = 0; s.cursor < len(s.instructions); s.cursor++ {
		if handler, found := m.handlers[s.instructions[s.cursor]]; found {
			if err := handler(m, s); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m Machine) Execute(input io.Reader, output io.Writer) error {
	instructions, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}

	m.output = output
	m.handlers = handlersSet

	return m.run(&scope{instructions: instructions})
}
