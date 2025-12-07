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

	startTime := time.Now()
	totalSplits := int64(0)

	// find the index of `S` in the first line
	sIndex := strings.Index(lines[0], "S")
	// build a slice of bools indicating tachyon beams (T) in that column, and initalize to false
	var tachyonBeams []int = make([]int, len(lines[0]))
	tachyonBeams[sIndex] = 1

	// process each line after the first one, updating tachyonBeams slice to carry the beams down,
	// splitting aroung `^` characters
	for _, line := range lines[1:] {
		for j, char := range line {
			currentBeams := tachyonBeams[j]
			if currentBeams == 0 {
				continue
			}
			if currentBeams > 0 && char == '^' {
				// split the beam
				if j > 0 {
					tachyonBeams[j-1] += currentBeams
				}
				if j < len(line)-1 {
					tachyonBeams[j+1] += currentBeams
				}
				tachyonBeams[j] = 0
				totalSplits++
			}
		}
	}
	fmt.Printf("Total splits: %d\n", totalSplits)
	fmt.Printf("Total active timelines: %d\n", func() int64 {
		var sum int64 = 0
		for _, v := range tachyonBeams {
			sum += int64(v)
		}
		return sum
	}())
	fmt.Printf("Execution time: %s\n", time.Since(startTime))
}
