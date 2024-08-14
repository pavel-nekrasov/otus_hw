package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type ParseState struct {
	builder strings.Builder
	Cur     rune
	Prev    rune
	Escaped bool
}

func Unpack(input string) (string, error) {
	var state ParseState

	for _, state.Cur = range input {
		if state.Escaped {
			if !state.isBackSlash() && !state.isDigit() {
				return "", ErrInvalidString
			}
			state.Escaped = false
			state.push()
			continue
		}

		switch {
		case state.isBackSlash():
			state.output()
			state.Escaped = true
		case state.isEmpty() && state.isDigit():
			return "", ErrInvalidString
		case state.isEmpty():
			state.push()
		case state.isDigit():
			state.outputN()
		default:
			state.output()
			state.push()
		}
	}

	state.output()

	return state.string(), nil
}

func (s *ParseState) isDigit() bool {
	return s.Cur >= 48 && s.Cur <= 57
}

func (s *ParseState) isBackSlash() bool {
	return s.Cur == 92
}

func (s *ParseState) toDigit() int {
	return int(s.Cur) - 48
}

func (s *ParseState) isEmpty() bool {
	return s.Prev == 0
}

func (s *ParseState) push() {
	s.Prev = s.Cur
}

func (s *ParseState) pop() rune {
	r := s.Prev
	s.Prev = 0
	return r
}

func (s *ParseState) output() {
	if !s.isEmpty() {
		s.builder.WriteString(string(s.pop()))
	}
}

func (s *ParseState) outputN() {
	if !s.isEmpty() {
		s.builder.WriteString(strings.Repeat(string(s.pop()), s.toDigit()))
	}
}

func (s *ParseState) string() string {
	return s.builder.String()
}
