package main

import (
	"github.com/Laky-64/gologging"
	"os"
)

func main() {
	file, _ := os.OpenFile("foxy_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer func() {
		_ = file.Close()
	}()
	gologging.SetOutput(file)
	gologging.Info("Hello, World!")
}
