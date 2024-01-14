package main

import (
	"os"
)

func someOtherFunc() {
	os.Exit(1)
}

func main() {
	//fmt.Println("Hello world")

	os.Exit(1) // want "call os.Exit in main func"
	someOtherFunc()

	//fmt.Println("Unreachable code")
}
