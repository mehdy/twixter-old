package entities

// Logger defines logging behaviour.
type Logger interface {
	// As sets the level of log. From 0 to 4 for Debug, Info, Warn, Error and Fatal respectively.
	As(level int8) Logger
	// WithField adds a pair of key/value to the log entry.
	WithField(key string, value interface{}) Logger
	// Logf emits the log message with formatting capabilities.
	Logf(msg string, args ...interface{})
}
