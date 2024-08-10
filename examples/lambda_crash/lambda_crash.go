package main

import "github.com/Laky-64/gologging"

func main() {
	testLambda := func() {
		gologging.Fatal("I'm a lambda, and you can't stop me from crashing!")
	}
	testLambda()
}
