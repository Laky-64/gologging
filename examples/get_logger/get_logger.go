package main

import "github.com/Laky-64/gologging"

func main() {
	logger := gologging.GetLogger("test")
	logger.Fatal("Logging from logger")
}
