package logger

import "testing"

func TestNewInternalLogger(t *testing.T) {
	logger := newInternalLogger(LevelDebug)

	logger.Debug("测试D")
	logger.Info("测试I")
	logger.Warn("测试W")
	logger.Error("测试E")
	logger.Fatal("测试F")
}
