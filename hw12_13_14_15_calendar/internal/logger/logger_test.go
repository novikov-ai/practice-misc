package logger

import (
	"fmt"
	"strings"
	"testing"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/configs"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	testCases := []struct {
		logLevel string
		output   []string
	}{
		{logLevel: "DEBUG", output: []string{"DEBUG", "INFO", "WARN", "ERROR"}},
		{logLevel: "INFO", output: []string{"INFO", "WARN", "ERROR"}},
		{logLevel: "WARN", output: []string{"WARN", "ERROR"}},
		{logLevel: "ERROR", output: []string{"ERROR"}},
	}

	for _, test := range testCases {
		test := test

		t.Run(fmt.Sprintf("Testing log-level: %s\n", test.logLevel), func(t *testing.T) {
			config := configs.Config{}
			config.Logger.Level = test.logLevel

			logger := New(config)

			w := CustomWriter{}
			logger.Writer = &w

			logged := make([]string, 0)

			logger.Debug("...")
			flushed := w.Flush()
			if flushed != "" {
				logged = append(logged, flushed)
			}

			logger.Info("...")
			flushed = w.Flush()
			if flushed != "" {
				logged = append(logged, flushed)
			}

			logger.Warn("...")
			flushed = w.Flush()
			if flushed != "" {
				logged = append(logged, flushed)
			}

			logger.Error("...")
			flushed = w.Flush()
			if flushed != "" {
				logged = append(logged, flushed)
			}

			require.Equal(t, len(test.output), len(logged))

			for i := 0; i < len(logged); i++ {
				require.True(t, strings.Contains(logged[i], test.output[i]))
			}
		})
	}
}

type CustomWriter struct {
	content strings.Builder
}

func (w *CustomWriter) Write(p []byte) (n int, err error) {
	return w.content.Write(p)
}

func (w *CustomWriter) Flush() string {
	defer w.content.Reset()
	return w.content.String()
}
