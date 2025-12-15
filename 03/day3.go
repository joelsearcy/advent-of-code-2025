package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	startTime := time.Now()

	// split the input data into lines
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	totalJoltage := 0
	batteryBankSize := 12

	maxJoltage := make([]byte, 0, batteryBankSize)
	// find the ordered "set" of digits (not necessarily adjacent) that make the largest N-digit number possible in each line
	for _, line := range lines {
		maxJoltage = maxJoltage[:0]
		startIndex := 0
		for i := range batteryBankSize {
			endIndex := len(line) - (batteryBankSize - i) + 1
			if startIndex >= endIndex {
				break
			}
			maxDigitIndex, maxDigit := findIndexOfNextMaxDigit(line[:endIndex], startIndex)
			maxJoltage = append(maxJoltage, maxDigit)
			startIndex = maxDigitIndex + 1
		}
		joltage, err := strconv.Atoi(string(maxJoltage))
		if err != nil {
			panic(err)
		}
		totalJoltage += joltage
	}

	fmt.Printf("Total Joltage: %d\n", totalJoltage)
	fmt.Printf("Execution time: %s\n", time.Since(startTime))
}

func findIndexOfNextMaxDigit(line string, startIndex int) (int, byte) {
	maxDigit := line[startIndex]
	maxDigitIndex := startIndex
	for i := startIndex + 1; i < len(line); i++ {
		if line[i] > maxDigit {
			maxDigit = line[i]
			maxDigitIndex = i
		}
	}
	return maxDigitIndex, maxDigit
}
