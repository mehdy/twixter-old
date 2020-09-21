package logrus

import (
	"github.com/mehdy/twixter/pkg/entities"
	"github.com/sirupsen/logrus"
)

// Ensure Logger and Entry implement entities.Logger interface.
var (
	_ entities.Logger = &Logger{}
	_ entities.Logger = &Entry{}
)

type Logger struct {
	*logrus.Logger
}

type Entry struct {
	*logrus.Entry
}

func NewLogger(config entities.ConfigGetter) *Logger {
	logger := logrus.New()
	logger.SetLevel(toLogrusLevel(config.GetString("log.level")))

	return &Logger{Logger: logger}
}

func (e *Entry) As(level string) entities.Logger {
	e.Level = toLogrusLevel(level)

	return e
}

func (e *Entry) WithField(key string, value interface{}) entities.Logger {
	e.Entry.WithField(key, value)

	return e
}

func (e *Entry) WithError(err error) entities.Logger {
	e.Entry.WithError(err)

	return e
}

func (e *Entry) Logf(msg string, args ...interface{}) {
	e.Entry.Logf(e.Level, msg, args...)
}

func (l *Logger) As(level string) entities.Logger {
	e := &Entry{logrus.NewEntry(l.Logger)}
	e.Level = toLogrusLevel(level)

	return e
}

func (l *Logger) WithField(key string, value interface{}) entities.Logger {
	return &Entry{l.Logger.WithField(key, value)}
}

func (l *Logger) WithError(err error) entities.Logger {
	return &Entry{l.Logger.WithError(err)}
}

func (l *Logger) Logf(msg string, args ...interface{}) {
	l.Logger.Logf(l.Level, msg, args...)
}

func toLogrusLevel(level string) logrus.Level {
	switch level {
	case "D":
		return logrus.DebugLevel
	case "I":
		return logrus.InfoLevel
	case "W":
		return logrus.WarnLevel
	case "E":
		return logrus.ErrorLevel
	case "F":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
