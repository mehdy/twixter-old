package entities

// Logger defines logging behaviour.
type Logger interface {
	// As sets the level of log. "D", "I", "W", "E", "F" for Debug, Info, Warn, Error and Fatal respectively.
	As(level string) Logger
	// WithField adds a pair of key/value to the log entry.
	WithField(key string, value interface{}) Logger
	// WithError adds an error value to the log entry.
	WithError(err error) Logger
	// Logf emits the log message with formatting capabilities.
	Logf(msg string, args ...interface{})
}
