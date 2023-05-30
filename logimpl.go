package log

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// ReLogger is the implemention for logrus.
type ReLogger struct {
	// name is the name of logger that is published to log as a scope
	name string
	// loger is the instance of logrus logger
	logger *logrus.Entry
}

var DaprVersion = "unknown"

func newReLogger(name string) *ReLogger {
	newLogger := logrus.New()
	newLogger.SetOutput(os.Stdout)

	dl := &ReLogger{
		name: name,
		logger: newLogger.WithFields(logrus.Fields{
			logFieldScope: name,
			logFieldType:  LogTypeLog,
		}),
	}

	dl.EnableJSONOutput(defaultJSONOutput)

	return dl
}

// EnableJSONOutput enables JSON formatted output log.
func (l *ReLogger) EnableJSONOutput(enabled bool) {
	var formatter logrus.Formatter

	fieldMap := logrus.FieldMap{
		// If time field name is conflicted, logrus adds "fields." prefix.
		// So rename to unused field @time to avoid the confliction.
		logrus.FieldKeyTime:  logFieldTimeStamp,
		logrus.FieldKeyLevel: logFieldLevel,
		logrus.FieldKeyMsg:   logFieldMessage,
	}

	hostname, _ := os.Hostname()
	l.logger.Data = logrus.Fields{
		logFieldScope:    l.logger.Data[logFieldScope],
		logFieldType:     LogTypeLog,
		logFieldInstance: hostname,
		logFieldDaprVer:  DaprVersion,
	}

	if enabled {
		formatter = &logrus.JSONFormatter{ //nolint: exhaustruct
			TimestampFormat: time.RFC3339Nano,
			FieldMap:        fieldMap,
		}
	} else {
		formatter = &logrus.TextFormatter{ //nolint: exhaustruct
			TimestampFormat: time.RFC3339Nano,
			FieldMap:        fieldMap,
		}
	}

	l.logger.Logger.SetFormatter(formatter)
}

// SetAppID sets app_id field in the log. Default value is empty string.
func (l *ReLogger) SetAppID(id string) {
	l.logger = l.logger.WithField(logFieldAppID, id)
}

func toLogrusLevel(lvl LogLevel) logrus.Level {
	// ignore error because it will never happen
	l, _ := logrus.ParseLevel(string(lvl))
	return l
}

// SetOutputLevel sets log output level.
func (l *ReLogger) SetOutputLevel(outputLevel LogLevel) {
	l.logger.Logger.SetLevel(toLogrusLevel(outputLevel))
}

// IsOutputLevelEnabled returns true if the logger will output this LogLevel.
func (l *ReLogger) IsOutputLevelEnabled(level LogLevel) bool {
	return l.logger.Logger.IsLevelEnabled(toLogrusLevel(level))
}

// SetOutput sets the destination for the logs.
func (l *ReLogger) SetOutput(dst io.Writer) {
	l.logger.Logger.SetOutput(dst)
}

// WithLogType specify the log_type field in log. Default value is LogTypeLog.
func (l *ReLogger) WithLogType(logType string) Logger {
	return &ReLogger{
		name:   l.name,
		logger: l.logger.WithField(logFieldType, logType),
	}
}

// WithFields returns a logger with the added structured fields.
func (l *ReLogger) WithFields(fields map[string]any) Logger {
	return &ReLogger{
		name:   l.name,
		logger: l.logger.WithFields(fields),
	}
}

// Info logs a message at level Info.
func (l *ReLogger) Info(args ...interface{}) {
	l.logger.Log(logrus.InfoLevel, args...)
}

// Infof logs a message at level Info.
func (l *ReLogger) Infof(format string, args ...interface{}) {
	l.logger.Logf(logrus.InfoLevel, format, args...)
}

// Debug logs a message at level Debug.
func (l *ReLogger) Debug(args ...interface{}) {
	l.logger.Log(logrus.DebugLevel, args...)
}

// Debugf logs a message at level Debug.
func (l *ReLogger) Debugf(format string, args ...interface{}) {
	l.logger.Logf(logrus.DebugLevel, format, args...)
}

// Warn logs a message at level Warn.
func (l *ReLogger) Warn(args ...interface{}) {
	l.logger.Log(logrus.WarnLevel, args...)
}

// Warnf logs a message at level Warn.
func (l *ReLogger) Warnf(format string, args ...interface{}) {
	l.logger.Logf(logrus.WarnLevel, format, args...)
}

// Error logs a message at level Error.
func (l *ReLogger) Error(args ...interface{}) {
	l.logger.Log(logrus.ErrorLevel, args...)
}

// Errorf logs a message at level Error.
func (l *ReLogger) Errorf(format string, args ...interface{}) {
	l.logger.Logf(logrus.ErrorLevel, format, args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1.
func (l *ReLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf logs a message at level Fatal then the process will exit with status set to 1.
func (l *ReLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
