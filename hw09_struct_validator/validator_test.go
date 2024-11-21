package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	InnerStruct struct {
		InnerVal1 string `validate:"len:7"`
	}

	OuterStruct1 struct {
		OuterVal1 string `validate:"len:5"`
		Inner     InnerStruct
		OuterVal2 string `validate:"len:5"`
	}

	OuterStruct2 struct {
		OuterVal1 string      `validate:"len:5"`
		Inner     InnerStruct `validate:"nested"`
		OuterVal2 string      `validate:"len:5"`
	}

	OuterStruct3 struct {
		OuterVal1 string       `validate:"len:5"`
		Inner     *InnerStruct `validate:"nested"`
		OuterVal2 string       `validate:"len:5"`
	}

	BadStruct1 struct {
		Va1 string `validate:"len:xxx"`
	}
	BadStruct2 struct {
		Va1 int `validate:"min:10|max:xxx"`
	}
	BadStruct3 struct {
		Va1 int `validate:"len:11"`
	}
	BadStruct4 struct {
		Va1 int `validate:"in:100,200,xxx"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          User{ID: "12345", Name: "somename", Age: 15, Email: "sssddd", Role: "stuff"},
			expectedErr: fmt.Errorf("ID: wrong len 5, expected 36\nAge: wrong value 15, should greater or equal to 18\nEmail: wrong value \"sssddd\""),
		},
		{
			in:          User{ID: "12345", Name: "somename", Age: 125, Email: "user@example.com", Role: "admin"},
			expectedErr: fmt.Errorf("ID: wrong len 5, expected 36\nAge: wrong value 125, should less or equal to 50"),
		},
		{
			in:          User{ID: "12345", Name: "somename", Age: 45, Email: "user@example.com", Role: "unknown"},
			expectedErr: fmt.Errorf("ID: wrong len 5, expected 36\nRole: wrong value \"unknown\", should be one of [admin stuff]"),
		},
		{
			in:          User{ID: "012345678901234567890123456789012345", Name: "somename", Age: 45, Email: "user@example.com", Role: "stuff"},
			expectedErr: nil,
		},
		{
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		{
			in:          App{Version: "1234"},
			expectedErr: fmt.Errorf("Version: wrong len 4, expected 5"),
		},
		{
			in:          Token{Header: []byte{1, 2, 3}},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 500},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 100, Body: "somevalue"},
			expectedErr: fmt.Errorf("Code: wrong value 100, should be one of [200 404 500]"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func TestNestedWithPointersValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          OuterStruct1{OuterVal1: "xxx", OuterVal2: "yyyyyyy", Inner: InnerStruct{InnerVal1: "zzzz"}},
			expectedErr: fmt.Errorf("OuterVal1: wrong len 3, expected 5\nOuterVal2: wrong len 7, expected 5"),
		},
		{
			in:          &OuterStruct1{OuterVal1: "xxx", OuterVal2: "yyyyyyy", Inner: InnerStruct{InnerVal1: "zzzz"}},
			expectedErr: fmt.Errorf("OuterVal1: wrong len 3, expected 5\nOuterVal2: wrong len 7, expected 5"),
		},
		{
			in:          OuterStruct2{OuterVal1: "xxx", OuterVal2: "yyyyyyy", Inner: InnerStruct{InnerVal1: "zzzz"}},
			expectedErr: fmt.Errorf("OuterVal1: wrong len 3, expected 5\nInner.InnerVal1: wrong len 4, expected 7\nOuterVal2: wrong len 7, expected 5"),
		},
		{
			in:          OuterStruct3{OuterVal1: "xxx", OuterVal2: "yyyyyyy", Inner: &InnerStruct{InnerVal1: "zzzz"}},
			expectedErr: fmt.Errorf("OuterVal1: wrong len 3, expected 5\nInner.InnerVal1: wrong len 4, expected 7\nOuterVal2: wrong len 7, expected 5"),
		},
		{
			in:          OuterStruct3{OuterVal1: "xxxxx", OuterVal2: "yyyyy", Inner: &InnerStruct{InnerVal1: "zzzzzzz"}},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func TestBadNotationSyntax(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          BadStruct1{Va1: "xxx"},
			expectedErr: fmt.Errorf("Va1: wrong notation: len \"xxx\" is not an integer"),
		},
		{
			in:          BadStruct2{Va1: 15},
			expectedErr: fmt.Errorf("Va1: wrong notation: max \"xxx\" is not an integer"),
		},
		{
			in:          BadStruct3{Va1: 15},
			expectedErr: fmt.Errorf("Va1: wrong notation: unknown rule \"len\""),
		},
		{
			in:          BadStruct4{Va1: 100},
			expectedErr: fmt.Errorf("Va1: wrong notation: \"xxx\" is not an integer"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}
