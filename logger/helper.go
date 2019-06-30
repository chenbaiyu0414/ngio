package logger

var defaultLogger Logger = newInternalLogger(LevelDebug)

func DefaultLogger() Logger {
	return defaultLogger
}

func SetLogger(newLogger Logger) {
	if newLogger == nil {
		return
	}

	defaultLogger = newLogger
}
