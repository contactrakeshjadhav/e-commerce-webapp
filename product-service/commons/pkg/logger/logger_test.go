package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type LogMessage struct {
	Level   string
	Message string
}

type LogMessageWithContext struct {
	Level   string
	Message string
	Field1  string
	Field2  string
}

func GetLogger() *AppLogger {
	return NewLogger(LogMetadata{
		"version": "1.0",
	})
}

func TestCanFormatMessage(t *testing.T) {
	appLogger := GetLogger()
	appLogger.SetLevel(DEBUG)

	msg := "%d plus %d is %d"
	a := 2
	b := 5
	sum := a + b

	tests := []struct {
		exec  func(format string, v ...interface{})
		level string
	}{
		{
			exec:  appLogger.Errorf,
			level: "error",
		},
		{
			exec:  appLogger.Warnf,
			level: "warn",
		},
		{
			exec:  appLogger.Infof,
			level: "info",
		},
		{
			exec:  appLogger.Debugf,
			level: "debug",
		},
	}

	for _, test := range tests {
		t.Run(test.level, func(t *testing.T) {
			var buf bytes.Buffer
			log := zerolog.New(&buf).With().Logger()
			appLogger.logger = &log

			test.exec(msg, a, b, sum)

			var message LogMessage
			json.Unmarshal(buf.Bytes(), &message)
			require.Equal(t, fmt.Sprintf(msg, a, b, sum), message.Message)
			require.Equal(t, test.level, message.Level)
		})
	}
}

func TestCanAddContextByIncludingNewFields(t *testing.T) {
	appLogger := GetLogger()
	appLogger.SetLevel(DEBUG)

	var buf bytes.Buffer
	log := zerolog.New(&buf).With().Logger()
	appLogger.logger = &log

	msg := "A log message with multiple values"
	val1 := "value1"
	val2 := "value2"

	appLogger.Debug(msg, LogContext{
		"field1": val1,
		"field2": val2,
	})

	var message LogMessageWithContext
	json.Unmarshal(buf.Bytes(), &message)

	require.Equal(t, msg, message.Message)
	require.Equal(t, "debug", message.Level)
	require.Equal(t, val1, message.Field1)
	require.Equal(t, val2, message.Field2)
}

func TestCanLogError(t *testing.T) {
	appLogger := GetLogger()
	appLogger.SetLevel(ERROR)

	var buf bytes.Buffer
	log := zerolog.New(&buf).With().Logger()
	appLogger.logger = &log

	msg := "An error message"
	val := "value"

	appLogger.Error(msg, LogContext{
		"field1": val,
	})

	var message LogMessageWithContext
	json.Unmarshal(buf.Bytes(), &message)

	require.Equal(t, msg, message.Message)
	require.Equal(t, "error", message.Level)
	require.Equal(t, val, message.Field1)
}

func TestCanLogWarning(t *testing.T) {
	appLogger := GetLogger()
	appLogger.SetLevel(INFO)

	var buf bytes.Buffer
	log := zerolog.New(&buf).With().Logger()
	appLogger.logger = &log

	msg := "A warning message"
	val := "value"

	appLogger.Warn(msg, LogContext{
		"field1": val,
	})

	var message LogMessageWithContext
	json.Unmarshal(buf.Bytes(), &message)

	require.Equal(t, msg, message.Message)
	require.Equal(t, "warn", message.Level)
	require.Equal(t, val, message.Field1)
}

func TestCanLogInfo(t *testing.T) {
	appLogger := GetLogger()
	appLogger.SetLevel(INFO)

	var buf bytes.Buffer
	log := zerolog.New(&buf).With().Logger()
	appLogger.logger = &log

	msg := "An info message"
	val := "value"

	appLogger.Info(msg, LogContext{
		"field1": val,
	})

	var message LogMessageWithContext
	json.Unmarshal(buf.Bytes(), &message)

	require.Equal(t, msg, message.Message)
	require.Equal(t, "info", message.Level)
	require.Equal(t, val, message.Field1)
}

func TestCanTurnOffLoggerLevel(t *testing.T) {
	appLogger := GetLogger()
	appLogger.SetLevel(ERROR)

	var buf bytes.Buffer
	log := zerolog.New(&buf).With().Logger()
	appLogger.logger = &log

	appLogger.Warnf("Resource temporarily unavailable: %s", errors.New("fork: Unable to fork new process"))

	require.Empty(t, buf.String())
}

func TestCanGetLevelFromStr(t *testing.T) {
	require.Equal(t, DEBUG, StrToLevel("DEBUG"))
	require.Equal(t, INFO, StrToLevel("INFO"))
	require.Equal(t, WARN, StrToLevel("WARN"))
	require.Equal(t, ERROR, StrToLevel("ERROR"))
}

func TestLogFatal(t *testing.T) {
	appLogger := GetLogger()

	if os.Getenv("RUN_FATAL") == "1" {
		appLogger.Fatalf("cannot recover from this error")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLogFatal")
	cmd.Env = append(os.Environ(), "RUN_FATAL=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Want exit status 1")
}

type TestStruct struct {
	ID          string
	Names       []string
	Description string
	Age         int
}

func TestLogDTO(t *testing.T) {
	appLogger := GetLogger()
	appLogger.SetLevel(DEBUG)

	var buf bytes.Buffer
	log := zerolog.New(&buf).With().Logger()
	appLogger.logger = &log

	test := TestStruct{
		ID:          "234234-234-234-252-asdf9wefs3",
		Names:       []string{"Test", "Struct"},
		Description: "This is a test",
		Age:         32,
	}
	testOutput := "commons/pkg/logger.TestStruct: " +
		"{ID:234234-234-234-252-asdf9wefs3 Names:[Test Struct] Description:This is a test Age:32}"

	appLogger.LogDTO(test)

	var message LogMessageWithContext
	json.Unmarshal(buf.Bytes(), &message)

	require.Equal(t, testOutput, message.Message)
	require.Equal(t, "info", message.Level)
}
