package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	// split the input data into lines
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	totalRemoved := 0

	startTime := time.Now()

	for {
		currentPassRemoved := 0

		for i := range lines {
			removed := removeAccessibleCells(lines, i)
			currentPassRemoved += removed
		}
		if currentPassRemoved == 0 {
			break
		}
		totalRemoved += currentPassRemoved
	}

	fmt.Printf("Total accessable rolls: %d\n", totalRemoved)
	fmt.Printf("Execution time: %s\n", time.Since(startTime))
}

var adjacentOffsets = []struct{ rowOffset, colOffset int }{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func removeAccessibleCells(lines []string, currentRowIndex int) int {
	total := 0
	currentRow := []byte(lines[currentRowIndex])
	rowSize := len(currentRow)
	rowSetSize := len(lines)

	for colIndex := 0; colIndex < rowSize; colIndex++ {
		if currentRow[colIndex] != '@' {
			continue
		}

		adjacentAtCount := 0
		for _, offset := range adjacentOffsets {
			adjacentRowIndex := currentRowIndex + offset.rowOffset
			adjacentColIndex := colIndex + offset.colOffset

			if adjacentRowIndex < 0 || adjacentRowIndex >= rowSetSize ||
				adjacentColIndex < 0 || adjacentColIndex >= rowSize {
				continue
			}
			if lines[adjacentRowIndex][adjacentColIndex] == '@' {
				adjacentAtCount++
			}
		}

		if adjacentAtCount < 4 {
			total++
			currentRow[colIndex] = '.'
		}
	}

	lines[currentRowIndex] = string(currentRow)
	return total
}
