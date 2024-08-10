package main

import "github.com/Laky-64/gologging"

func main() {
	gologging.Debug("Setting up kernel of the protogen")
	gologging.Info("Protogen started")
	gologging.Warn("Protogen is warming up")
	gologging.Error("Protogen failed to start")
	gologging.Fatal("Protogen crashed")
}
