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
		log.New(debugWriter, prefixWithColor("[DEBUG]\t", colorGreen), log.LstdFlags|log.Lshortfile),
		log.New(infoWriter, prefixWithColor("[INFO]\t", colorBlue), log.LstdFlags|log.Lshortfile),
		log.New(warnWriter, prefixWithColor("[WARN]\t", colorYellow), log.LstdFlags|log.Lshortfile),
		log.New(errorWriter, prefixWithColor("[ERROR]\t", colorRed), log.LstdFlags|log.Lshortfile),
		log.New(fatalWriter, prefixWithColor("[FATAL]\t", colorMagenta), log.LstdFlags|log.Lshortfile),
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
	il.loggers[LevelDebug].Output(2, fmt.Sprint(v...))
}

func (il *internalLogger) Debugf(format string, v ...interface{}) {
	il.loggers[LevelDebug].Output(2, fmt.Sprintf(format, v...))
}

func (il *internalLogger) Info(v ...interface{}) {
	il.loggers[LevelInfo].Output(2, fmt.Sprint(v...))
}

func (il *internalLogger) Infof(format string, v ...interface{}) {
	il.loggers[LevelInfo].Output(2, fmt.Sprintf(format, v...))
}

func (il *internalLogger) Warn(v ...interface{}) {
	il.loggers[LevelWarn].Output(2, fmt.Sprint(v...))
}

func (il *internalLogger) Warnf(format string, v ...interface{}) {
	il.loggers[LevelWarn].Output(2, fmt.Sprintf(format, v...))
}

func (il *internalLogger) Error(v ...interface{}) {
	il.loggers[LevelError].Output(2, fmt.Sprint(v...))
}

func (il *internalLogger) Errorf(format string, v ...interface{}) {
	il.loggers[LevelError].Output(2, fmt.Sprintf(format, v...))
}

func (il *internalLogger) Fatal(v ...interface{}) {
	il.loggers[LevelFatal].Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func (il *internalLogger) Fatalf(format string, v ...interface{}) {
	il.loggers[LevelFatal].Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func prefixWithColor(prefix string, color int) string {
	return fmt.Sprintf("%c[%sm%s%c[0m", 0x1B, strconv.Itoa(color), prefix, 0x1B)
}
