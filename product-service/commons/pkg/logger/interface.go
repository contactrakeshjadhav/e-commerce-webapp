package logger

import "context"

type Logger interface {
	// SetLevel will set the desired logging level.
	SetLevel(level LogLevel)

	// SetLevelFromDefaultEnvVar Will get the value from configuration setting LOG_LEVEL
	// and set it as level. If it is empty, the default level will be set, which is INFO
	SetLevelFromDefaultEnvVar()

	// Error writes a log with ERROR level. It accepts a map of strings as context data that can hold any additional information.
	Error(msg string, context LogContext)
	// Warn writes a log with WARN level. It accepts a map of strings as context data that can hold any additional information.
	Warn(msg string, context LogContext)
	// Info writes a log with INFO level. It accepts a map of strings as context data that can hold any additional information.
	Info(msg string, context LogContext)
	// Debug writes a log with DEBUG level. It accepts a map of strings as context data that can hold any additional information.
	Debug(msg string, context LogContext)

	// Errorf formats message according to a format specifier and writes a log with ERROR level.
	Errorf(format string, v ...interface{})
	// Warnf formats message according to a format specifier and writes a log with WARN level.
	Warnf(format string, v ...interface{})
	// Infof formats message according to a format specifier and writes a log with INFO level.
	Infof(format string, v ...interface{})
	// Debugf formats message according to a format specifier and writes a log with DEBUG level.
	Debugf(format string, v ...interface{})

	// Fatalf will log the message and then terminate the program.
	Fatalf(format string, v ...interface{})

	// LogDTO will log the struct with INFO level.
	LogDTO(v interface{})

	// WithStr will add a key:value pair of strings to the logger context
	WithStr(key, value string) Logger

	// WithInt will a key and a int value to the logger context
	WithInt(key string, value int) Logger

	// WithReqId will add a key:value pair of strings that consist of ReqIDKey:ReqID
	WithReqID(context context.Context) Logger
}
