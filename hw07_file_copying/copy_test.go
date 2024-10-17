package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
)

func TestCopyInvalidCases(t *testing.T) {
	t.Run("source path does not exist", func(t *testing.T) {
		err := Copy("./testdata/not_existent.txt", "./", 0, 0)
		require.Truef(t, errors.Is(err, ErrSourceFileNotFound), "actual err - %v", err)
	})

	t.Run("target path does not exist", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "./testdata_wrong/output.txt", 0, 0)
		require.Truef(t, errors.Is(err, ErrTargetCannotBeCreated), "actual err - %v", err)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "/tmp/output.txt", 100000, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})
}

func TestCopy(t *testing.T) {
	tests := []struct {
		from   string
		offset int64
		limit  int64

		example string
	}{
		{
			from:    "./testdata/input.txt",
			offset:  0,
			limit:   0,
			example: "./testdata/out_offset0_limit0.txt",
		},
		{
			from:    "./testdata/input.txt",
			offset:  0,
			limit:   10,
			example: "./testdata/out_offset0_limit10.txt",
		},
		{
			from:    "./testdata/input.txt",
			offset:  0,
			limit:   1000,
			example: "./testdata/out_offset0_limit1000.txt",
		},
		{
			from:    "./testdata/input.txt",
			offset:  0,
			limit:   10000,
			example: "./testdata/out_offset0_limit10000.txt",
		},
		{
			from:    "./testdata/input.txt",
			offset:  100,
			limit:   1000,
			example: "./testdata/out_offset100_limit1000.txt",
		},
		{
			from:    "./testdata/input.txt",
			offset:  6000,
			limit:   1000,
			example: "./testdata/out_offset6000_limit1000.txt",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.example, func(t *testing.T) {
			fOut, err := os.CreateTemp("", "out.*.txt")
			if err != nil {
				panic(err)
			}
			err = Copy(tc.from, fOut.Name(), tc.offset, tc.limit)
			defer os.Remove(fOut.Name())

			require.NoError(t, err)

			cmp := equalfile.New(nil, equalfile.Options{})
			result, _ := cmp.CompareFile(fOut.Name(), tc.example)
			require.True(t, result)
		})
	}
}
