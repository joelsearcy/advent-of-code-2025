package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	// Split the input data into ranges
	ranges := strings.Split(strings.TrimSpace(string(data)), ",")
	totals := 0

	for _, r := range ranges {
		// fmt.Printf("Range: %s\n", r)
		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			fmt.Printf("Invalid range: %s\n", r)
			continue
		}
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		// fmt.Printf("Start: %d, End: %d\n", start, end)

		for i := start; i <= end; i++ {
			n := fmt.Sprintf("%d", i)
			// pattern: even number of digits and first half == second half
			if len(n)%2 == 0 && n[0:len(n)/2] == n[len(n)/2:] {
				// fmt.Printf("  Found invalid ID: %d\n", i)
				totals += i
				continue
			}

			// or, any pattern of 1 or more digits repeated 2 or more times consecutively to form the entire number
			for l := 1; l <= len(n)/2; l++ {
				if len(n)%l == 0 {
					repeats := len(n) / l
					pattern := n[0:l]
					valid := true
					for j := 1; j < repeats; j++ {
						if n[j*l:(j+1)*l] != pattern {
							valid = false
							break
						}
					}
					if valid {
						totals += i
						break
					}
				}
			}
		}
	}
	fmt.Printf("Total of all invalid IDs: %d\n", totals)
}
