package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	colorRed     = 31
	colorGreen   = 32
	colorYellow  = 33
	colorBlue    = 34
	colorMagenta = 35
)

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type internalLogger struct {
	loggers []*log.Logger
}

func NewInternalLogger(debugWriter, infoWriter, warnWriter, errorWriter, fatalWriter io.Writer) *internalLogger {
	loggers := []*log.Logger{
		log.New(debugWriter, prefixWithColor("[DEBUG] ", colorGreen), log.LstdFlags),
		log.New(infoWriter, prefixWithColor("[INFO] ", colorBlue), log.LstdFlags),
		log.New(warnWriter, prefixWithColor("[WARN] ", colorYellow), log.LstdFlags),
		log.New(errorWriter, prefixWithColor("[ERROR] ", colorRed), log.LstdFlags),
		log.New(fatalWriter, prefixWithColor("[FATAL] ", colorMagenta), log.LstdFlags),
	}

	return &internalLogger{loggers: loggers}
}

func newInternalLogger(level Level) *internalLogger {
	writers := [5]io.Writer{
		os.Stdout, // Debug
		os.Stdout, // Info
		os.Stdout, // Warn
		os.Stdout, // Error
		os.Stdout, // Fatal
	}

	for writerLevel := range writers {
		if writerLevel < int(level) {
			writers[writerLevel] = ioutil.Discard
		}
	}

	return NewInternalLogger(writers[LevelDebug], writers[LevelInfo], writers[LevelWarn], writers[LevelError], writers[LevelFatal])
}

func (il *internalLogger) Debug(v ...interface{}) {
	il.loggers[LevelDebug].Print(v...)
}

func (il *internalLogger) Debugf(format string, v ...interface{}) {
	il.loggers[LevelDebug].Printf(format, v...)
}

func (il *internalLogger) Info(v ...interface{}) {
	il.loggers[LevelInfo].Print(v...)
}

func (il *internalLogger) Infof(format string, v ...interface{}) {
	il.loggers[LevelInfo].Printf(format, v...)
}

func (il *internalLogger) Warn(v ...interface{}) {
	il.loggers[LevelWarn].Print(v...)
}

func (il *internalLogger) Warnf(format string, v ...interface{}) {
	il.loggers[LevelWarn].Printf(format, v...)
}

func (il *internalLogger) Error(v ...interface{}) {
	il.loggers[LevelError].Print(v...)
}

func (il *internalLogger) Errorf(format string, v ...interface{}) {
	il.loggers[LevelError].Printf(format, v...)
}

func (il *internalLogger) Fatal(v ...interface{}) {
	il.loggers[LevelFatal].Fatal(v...)
}

func (il *internalLogger) Fatalf(format string, v ...interface{}) {
	il.loggers[LevelFatal].Fatalf(format, v...)
}

func prefixWithColor(prefix string, color int) string {
	return fmt.Sprintf("%c[%sm%s%c[0m", 0x1B, strconv.Itoa(color), prefix, 0x1B)
}
