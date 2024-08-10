package gologging

import (
	"testing"
)

func TestStdLog(t *testing.T) {
	Debug("asasdf", 1)
	Info("asasdf", 1)
	Warn("asasdf", 1)
	Error("asasdf", 1)
	Fatal("asasdf", 1)
}
