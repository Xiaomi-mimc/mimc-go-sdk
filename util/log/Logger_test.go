package log

import (
	"testing"
)

func TestInfoLogger(test *testing.T) {
	SetLogLevel(WarnLevel)
	logger := GetLogger()
	logger.Info("%s.\n", "hello world")
	logger.Info("%s.\n", "hello world1")
	logger.Warn("%s.\n", "hello world2")
}
