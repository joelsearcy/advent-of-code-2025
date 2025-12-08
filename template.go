package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	data, err := os.ReadFile("sample.txt")
	if err != nil {
		panic(err)
	}
	// split the input data into lines
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	startTime := time.Now()
	total := int64(0)

	// start here

	fmt.Printf("Total: %d\n", total)
	fmt.Printf("Execution time: %s\n", time.Since(startTime))
}
