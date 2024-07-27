package logger

// LogMetadata adds each key-value pair to the logger context.
type LogMetadata map[string]interface{}

// LogContext represents context data for a message, key-value pairs with relevant information that helps to understand what happened
type LogContext map[interface{}]interface{}

// LogLevel represents a log level. It determines how important a given message is.
// We defined the followings levels: DEBUG, INFO, WARNING, and ERROR
type LogLevel int8

const LOG_LEVEL string = "LOG_LEVEL"

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)
