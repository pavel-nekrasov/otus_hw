package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	actualEnvs, err := ReadDir("./testdata/env")
	require.NoError(t, err)

	tests := []struct {
		key        string
		value      string
		needRemove bool
	}{
		{
			key:        "BAR",
			value:      "bar",
			needRemove: false,
		},
		{
			key:        "EMPTY",
			value:      "",
			needRemove: false,
		},
		{
			key:        "FOO",
			value:      "   foo\nwith new line",
			needRemove: false,
		},
		{
			key:        "HELLO",
			value:      "\"hello\"",
			needRemove: false,
		},
		{
			key:        "UNSET",
			value:      "",
			needRemove: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.key, func(t *testing.T) {
			actual, ok := actualEnvs[tc.key]
			require.True(t, ok)

			require.Equal(t, tc.needRemove, actual.NeedRemove)
			require.Equal(t, tc.value, actual.Value)
		})
	}
}

func TestBadCases(t *testing.T) {
	t.Run("Wrong env path", func(t *testing.T) {
		_, err := ReadDir("./testdata/env_wrong_path")
		require.Error(t, err)
	})

	t.Run("Wrong env file", func(t *testing.T) {
		_, err := ReadDir("./testdata/env_bad")
		require.Error(t, err, errWrongName)
	})
}
