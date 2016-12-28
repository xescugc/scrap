package main

import "fmt"

var (
	version = "1.0.1"
)

func printVersion() {
	fmt.Printf("The actual version of 'scrap' is: %v\n", version)
}
