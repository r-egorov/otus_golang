package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	testCases := []struct {
		name, level string
	}{
		{
			name:  "DEBUG LEVEL",
			level: "DEBUG",
		},
		{
			name:  "INFO LEVEL",
			level: "INFO",
		},
		{
			name:  "WARN LEVEL",
			level: "WARN",
		},
		{
			name:  "ERROR LEVEL",
			level: "ERROR",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			msgErr := "err message"
			msgWarn := "warn message"
			msgInfo := "info message"
			msgDebug := "debug message"

			l := New(out, tc.level)
			l.Error(msgErr)
			l.Warn(msgWarn)
			l.Info(msgInfo)
			l.Debug(msgDebug)

			errorPrefix := "[ERROR]"
			warnPrefix := "[WARN]"
			infoPrefix := "[INFO]"
			debugPrefix := "[DEBUG]"

			var expected string
			switch tc.level {
			case "ERROR":
				expected = errorPrefix + " " + msgErr + "\n"
			case "WARN":
				expected = errorPrefix + " " + msgErr + "\n" +
					warnPrefix + " " + msgWarn + "\n"
			case "INFO":
				expected = errorPrefix + " " + msgErr + "\n" +
					warnPrefix + " " + msgWarn + "\n" +
					infoPrefix + " " + msgInfo + "\n"
			case "DEBUG":
				expected = errorPrefix + " " + msgErr + "\n" +
					warnPrefix + " " + msgWarn + "\n" +
					infoPrefix + " " + msgInfo + "\n" +
					debugPrefix + " " + msgDebug + "\n"
			}

			require.Equal(t, expected, out.String())
		})
	}
}
