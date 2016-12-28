package main

import "fmt"

var (
	version = "1.0.3"
)

func printVersion() {
	fmt.Printf("The actual version of 'scrap' is: %v\n", version)
}
