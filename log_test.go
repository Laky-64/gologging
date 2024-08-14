package gologging

import (
	"testing"
	"time"
)

func TestStdLog(t *testing.T) {
	Debug("asasdf", 1)
	Info("asasdf", 1)
	Warn("asasdf", 1)
	Error("asasdf", 1)
	Fatal("asasdf", 1)
}

func TestBrokenGoRoutine(t *testing.T) {
	go Fatal("Error while sending message to Nick: Bad Request: message text is empty")
	time.Sleep(1 * time.Second)
}

func TestLambda(t *testing.T) {
	func() {
		Fatal("Error while sending message to Nick: Bad Request: message text is empty")
	}()
}
