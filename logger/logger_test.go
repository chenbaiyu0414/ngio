package logger

import (
	"testing"
)

func TestNewInternalLogger(t *testing.T) {
	l := newInternalLogger(LevelDebug)

	l.Debug("测试D")
	l.Info("测试I")
	l.Warn("测试W")
	l.Error("测试E")
	l.Fatal("测试F")
}
