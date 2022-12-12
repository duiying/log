package log

import (
	"testing"
)

var Log = NewLogger(DefaultConfig)

func TestLog(t *testing.T) {
	TestWith(t)
	Log.Info("infoMsg")
	TestTrace(t)
}

func TestTrace(t *testing.T) {
	Log.Error("errorMsg", "key1", "val1", "key2", 2)
}

func TestWith(t *testing.T) {
	Log.With("uid", 2).Info("withMsg", "key3", "val3")
}
