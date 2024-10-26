package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Test echo", func(t *testing.T) {
		tempOut, err := os.CreateTemp("", "out.*.txt")
		require.NoError(t, err)

		defer tempOut.Close()

		env, err := ReadDir("./testdata/my_env")
		require.NoError(t, err)

		rescueStdout := os.Stdout
		os.Stdout = tempOut

		retCode := RunCmd([]string{"/bin/bash", "./testdata/my_echo.sh"}, env)

		os.Stdout = rescueStdout
		tempOut.Seek(0, 0)
		output, _ := io.ReadAll(tempOut)

		require.Equal(t, 0, retCode)
		require.Equal(t, "XXX is (value1), YYY is (value2)\n", string(output))
	})
}
