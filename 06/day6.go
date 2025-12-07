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

	// split the input data into lines
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	// numbers are positive integers in vertical columns, read top-to-bottom right-to-left
	// a column of only spaces is a separator between groups of numbers

	startTime := time.Now()
	total := int64(0)

	numLines := len(lines)
	if numLines < 2 {
		panic("not enough lines in input")
	}

	// last line is the operators for each column
	operators := strings.Fields(lines[numLines-1])
	numGroups := len(operators)

	var totals []int64 = make([]int64, numGroups)

	// process each line except the last one
	currentGroup := 0
	j := 0

	numStr := make([]byte, numLines-1)
	for currentGroup < numGroups {
		for j < len(lines[0]) {
			// reset numStr for each column
			numStr = numStr[:0]
			foundDigit := false

			// collect digits per column to form numbers, then apply the operator
			for i := 0; i < numLines-1; i++ {
				digit := lines[i][j]
				if digit != ' ' {
					numStr = append(numStr, digit)
				}
			}

			if len(numStr) > 0 {
				foundDigit = true
				num, err := strconv.ParseInt(string(numStr), 10, 64)
				if err != nil {
					panic(err)
				}
				// fmt.Printf("Group %d: applying %s to %d\n", currentGroup, operators[currentGroup], num)

				switch operators[currentGroup] {
				case "+":
					totals[currentGroup] += num
				case "*":
					if totals[currentGroup] == 0 {
						totals[currentGroup] = 1
					}
					totals[currentGroup] *= num
				default:
					panic("unknown operator: " + operators[currentGroup])
				}
			}
			j++

			if !foundDigit {
				break
			}
		}

		currentGroup++
	}

	// sum up the totals
	for _, colTotal := range totals {
		total += colTotal
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Total: %d\n", total)
	fmt.Printf("Elapsed time: %s\n", elapsed)
}
