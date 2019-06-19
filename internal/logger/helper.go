package logger

var l = newInternalLogger(LevelDebug)

func Debugf(format string, v ...interface{}) {
	l.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	l.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	l.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	l.Errorf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	l.Fatalf(format, v...)
}
