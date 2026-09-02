package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		fmt.Println("Initializing clamav-service configuration...")
		os.Exit(0)
	}

	fmt.Println("Starting clamav-service...")
}
