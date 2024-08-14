package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type ParseState struct {
	builder strings.Builder
	Cur     rune
	prev    rune
	Escaped bool
}

func Unpack(input string) (string, error) {
	var state ParseState

	for _, state.Cur = range input {
		switch {
		case state.Escaped && !state.IsBackSlash() && !state.IsDigit():
			return "", ErrInvalidString
		case state.Escaped:
			state.Escaped = false
			state.Push()
		case state.IsBackSlash():
			state.Output()
			state.Escaped = true
		case state.IsEmpty() && state.IsDigit():
			return "", ErrInvalidString
		case state.IsEmpty():
			state.Push()
		case state.IsDigit():
			state.OutputN()
		default:
			state.Output()
			state.Push()
		}
	}

	state.Output()

	return state.String(), nil
}

func (s *ParseState) IsDigit() bool {
	return s.Cur >= 48 && s.Cur <= 57
}

func (s *ParseState) IsBackSlash() bool {
	return s.Cur == 92
}

func (s *ParseState) Digit() int {
	return int(s.Cur) - 48
}

func (s *ParseState) IsEmpty() bool {
	return s.prev == 0
}

func (s *ParseState) Push() {
	s.prev = s.Cur
}

func (s *ParseState) Pop() rune {
	r := s.prev
	s.prev = 0
	return r
}

func (s *ParseState) Output() {
	if !s.IsEmpty() {
		s.builder.WriteString(string(s.Pop()))
	}
}

func (s *ParseState) OutputN() {
	if !s.IsEmpty() && s.IsDigit() {
		s.builder.WriteString(strings.Repeat(string(s.Pop()), s.Digit()))
	}
}

func (s *ParseState) String() string {
	return s.builder.String()
}
