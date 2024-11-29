package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggingCases(t *testing.T) {
	tests := []struct {
		typeMsg  string
		msg      string
		expected string
	}{
		{
			typeMsg:  "info",
			msg:      "test info msg",
			expected: "\"msg\":\"test info msg\"",
		},
		{
			typeMsg:  "warn",
			msg:      "test warn msg",
			expected: "\"msg\":\"test warn msg\"",
		},
		{
			typeMsg:  "debug",
			msg:      "test debug msg",
			expected: "\"msg\":\"test debug msg\"",
		},
		{
			typeMsg:  "error",
			msg:      "test error msg",
			expected: "\"msg\":\"test error msg\"",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.typeMsg, func(t *testing.T) {
			fOut, err := os.CreateTemp("", "logger*.txt")
			require.NoError(t, err)
			defer os.Remove(fOut.Name())

			l := New("debug", fOut.Name())

			switch tc.typeMsg {
			case "info":
				l.Info(tc.msg)
			case "warn":
				l.Warn(tc.msg)
			case "error":
				l.Error(tc.msg)
			case "debug":
				l.Debug(tc.msg)
			}
			output, err := os.ReadFile(fOut.Name())
			require.NoError(t, err)
			strOutput := string(output)
			strings.Contains(strOutput, tc.expected)
			require.Contains(t, strOutput, tc.expected)
		})
	}
}
