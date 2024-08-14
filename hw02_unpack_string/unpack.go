package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type ParseState struct {
	Cur     rune
	prev    rune
	Escaped bool
}

func Unpack(input string) (string, error) {
	var builder strings.Builder
	var state ParseState

	for _, state.Cur = range input {
		switch {
		case state.Escaped && !state.IsBackSlash() && !state.IsDigit():
			return "", ErrInvalidString
		case state.Escaped:
			state.Escaped = false
			state.Push()
		case state.IsBackSlash():
			if !state.IsEmpty() {
				builder.WriteString(string(state.Pop()))
			}
			state.Escaped = true
		case state.IsEmpty() && state.IsDigit():
			return "", ErrInvalidString
		case state.IsEmpty():
			state.Push()
		case state.IsDigit():
			builder.WriteString(strings.Repeat(string(state.Pop()), state.Digit()))
		default:
			builder.WriteString(string(state.Pop()))
			state.Push()
		}
	}

	if state.Escaped {
		return "", ErrInvalidString
	}

	if !state.IsEmpty() {
		builder.WriteString(string(state.Pop()))
	}

	return builder.String(), nil
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
