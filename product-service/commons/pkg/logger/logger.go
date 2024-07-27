package logger

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/model"

	"github.com/rs/zerolog"
)

// mapLevels will map our defined levels to zerolog levels
var mapLevels = map[LogLevel]zerolog.Level{
	DEBUG: zerolog.DebugLevel,
	INFO:  zerolog.InfoLevel,
	WARN:  zerolog.WarnLevel,
	ERROR: zerolog.ErrorLevel,
}

// mapStrLevels will map a string to eacho logger level
var mapStrLevels = map[string]LogLevel{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
}

type AppLogger struct {
	logger *zerolog.Logger
	level  LogLevel
}

// StrToLevel accepts a string and returns the corresponding LogLevel
func GetDefaultLevel() LogLevel {
	return INFO
}

// StrToLevel accepts a string and returns the corresponding LogLevel
func StrToLevel(str string) LogLevel {
	if l, ok := mapStrLevels[str]; ok {
		return l
	}
	return GetDefaultLevel()
}

func NewLogger(metadata LogMetadata) *AppLogger {
	builder := zerolog.New(os.Stdout).With().Timestamp()
	var logger zerolog.Logger

	for field, value := range metadata {
		builder = builder.Str(field, value.(string))
	}
	logger = builder.Logger()
	return &AppLogger{
		logger: &logger,
	}
}

func (l *AppLogger) SetLevel(level LogLevel) {
	l.level = level
	zerolog.SetGlobalLevel(mapLevels[level])
}

func (l *AppLogger) log(level LogLevel, msg string, context LogContext) {
	var builder *zerolog.Event
	switch level {
	case ERROR:
		builder = l.logger.Error()
	case WARN:
		builder = l.logger.Warn()
	case INFO:
		builder = l.logger.Info()
	case DEBUG:
		builder = l.logger.Debug()
	}

	for f, v := range context {
		builder = builder.Str(f.(string), v.(string))
	}

	builder.Msg(msg)
}

func (l *AppLogger) Error(msg string, context LogContext) {
	l.log(ERROR, msg, context)
}

func (l *AppLogger) Warn(msg string, context LogContext) {
	l.log(WARN, msg, context)
}

func (l *AppLogger) Info(msg string, context LogContext) {
	l.log(INFO, msg, context)
}

func (l *AppLogger) Debug(msg string, context LogContext) {
	l.log(DEBUG, msg, context)
}

func (l *AppLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Caller(1).Msgf(format, v...)
}

func (l *AppLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Caller(1).Msgf(format, v...)
}

func (l *AppLogger) Infof(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}

func (l *AppLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

func (l *AppLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Caller(1).Msgf(format, v...)
}

func (l *AppLogger) LogDTO(v interface{}) {
	t := reflect.TypeOf(v)
	structName := t.PkgPath() + "." + t.Name()
	msg := fmt.Sprintf("%s: %+v", structName, v)
	l.logger.Info().Msg(msg)
}

func (l *AppLogger) WithStr(key, value string) Logger {
	builder := l.logger.With().Str(key, value).Logger()
	return &AppLogger{
		logger: &builder,
		level:  l.level,
	}
}

func (l *AppLogger) WithInt(key string, value int) Logger {
	builder := l.logger.With().Int(key, value).Logger()
	return &AppLogger{
		logger: &builder,
		level:  l.level,
	}
}

func (l *AppLogger) WithReqID(context context.Context) Logger {
	reqId := context.Value(model.ReqIDKey)
	if reqId != nil {
		builder := l.logger.With().Str(model.ReqIDKey.String(), reqId.(string)).Logger()
		return &AppLogger{
			logger: &builder,
			level:  l.level,
		}
	}
	builder := l.logger.With().Str(model.ReqIDKey.String(), "not-found").Logger()
	return &AppLogger{
		logger: &builder,
		level:  l.level,
	}
}

func (l *AppLogger) SetLevelFromDefaultEnvVar() {
	l.SetLevel(StrToLevel(os.Getenv(LOG_LEVEL)))
}
