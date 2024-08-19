package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type parseState struct {
	Cur     rune
	prev    rune
	Escaped bool
}

func Unpack(input string) (string, error) {
	var builder strings.Builder
	var state parseState

	for _, state.Cur = range input {
		switch {
		case state.Escaped && !state.isBackSlash() && !state.isDigit():
			return "", ErrInvalidString
		case state.Escaped:
			state.Escaped = false
			state.push()
		case state.isBackSlash():
			if state.any() {
				builder.WriteString(string(state.pop()))
			}
			state.Escaped = true
		case state.isEmpty() && state.isDigit():
			return "", ErrInvalidString
		case state.isEmpty():
			state.push()
		case state.isDigit():
			builder.WriteString(strings.Repeat(string(state.pop()), state.digit()))
		default:
			builder.WriteString(string(state.pop()))
			state.push()
		}
	}

	if state.Escaped {
		return "", ErrInvalidString
	}

	if state.any() {
		builder.WriteString(string(state.pop()))
	}

	return builder.String(), nil
}

func (s *parseState) isDigit() bool {
	return s.Cur >= 48 && s.Cur <= 57 // проверка по ASCII коду цифр - ASCII(0) - 48, ASCII(1) - 29 .... ASCII(9) - 57
}

func (s *parseState) isBackSlash() bool {
	return s.Cur == 92 // 92 - ASCII код '\'
}

func (s *parseState) digit() int {
	return int(s.Cur) - 48 // чтобы получить значение цифры - отмнимаем ASCII код 0 - 48
}

func (s *parseState) isEmpty() bool {
	return s.prev == 0
}

func (s *parseState) any() bool {
	return !s.isEmpty()
}

func (s *parseState) push() {
	s.prev = s.Cur
}

func (s *parseState) pop() rune {
	r := s.prev
	s.prev = 0
	return r
}
