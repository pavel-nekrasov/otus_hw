package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigCases(t *testing.T) {
	tests := []struct {
		configFile string
		level      string
		output     string
	}{
		{
			configFile: "../../testdata/configs/config1.toml",
			level:      "INFO",
			output:     "stdout",
		},
		{
			configFile: "../../testdata/configs/config2.toml",
			level:      "DEBUG",
			output:     "/var/logs/calendar.log",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.configFile, func(t *testing.T) {
			if _, err := os.Stat(tc.configFile); err != nil {
				require.NoError(t, err)
			}
			c := New(tc.configFile)

			require.Equal(t, tc.level, c.Logger.Level)
			require.Equal(t, tc.output, c.Logger.Output)
		})
	}
}
