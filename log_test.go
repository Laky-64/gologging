package gologging

import (
	"testing"
	"time"
)

func TestStdLog(t *testing.T) {
	SetLevel(DebugLevel)
	Debug("asasdf", 1)
	Info("asasdf", 1)
	Warn("asasdf", 1)
	Error("asasdf", 1)
	Fatal("asasdf", 1)
}

func TestBrokenGoRoutine(t *testing.T) {
	go Fatal("TestBrokenGoRoutine")
	time.Sleep(1 * time.Second)
}

func TestLambda(t *testing.T) {
	func() {
		Fatal("TestLambda")
	}()
}

func TestGetLogger(t *testing.T) {
	SetLevel(DebugLevel)
	logger := GetLogger("test")
	logger.Debug("asasdf", 1)
	logger.Info("asasdf", 1)
	logger.Warn("asasdf", 1)
	logger.Error("asasdf", 1)
	logger.Fatal("asasdf", 1)
}
