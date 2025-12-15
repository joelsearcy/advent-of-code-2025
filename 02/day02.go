package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	usePatternGen := flag.Bool("pattern", false, "Use pattern generation algorithm instead of brute force")
	flag.Parse()

	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	ranges := parseRanges(string(data))
	startTime := time.Now()

	var total int
	var algorithm string

	if *usePatternGen {
		algorithm = "Pattern Generation"
		total = solveWithPatternGeneration(ranges)
	} else {
		algorithm = "Brute Force (Optimized)"
		total = solveWithBruteForce(ranges)
	}

	fmt.Printf("Algorithm: %s\n", algorithm)
	fmt.Printf("Total of all invalid IDs: %d\n", total)
	fmt.Printf("Execution time: %s\n", time.Since(startTime))
}

// Parse input data into Range structs
func parseRanges(data string) []Range {
	rangeStrs := strings.Split(strings.TrimSpace(data), ",")
	ranges := make([]Range, 0, len(rangeStrs))

	for _, r := range rangeStrs {
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
		ranges = append(ranges, Range{start, end})
	}

	return ranges
}

// Brute force approach: check every number in each range
func solveWithBruteForce(ranges []Range) int {
	total := 0
	buf := make([]byte, 0, 20)

	for _, r := range ranges {
		for i := r.start; i <= r.end; i++ {
			if hasRepeatingPattern(i, buf) {
				total += i
			}
		}
	}

	return total
}

// Pattern generation approach: generate only valid pattern numbers
func solveWithPatternGeneration(ranges []Range) int {
	maxNum := 0
	for _, r := range ranges {
		if r.end > maxNum {
			maxNum = r.end
		}
	}

	maxDigits := countDigits(maxNum)
	validNumbers := make(map[int]bool)

	// Generate all possible pattern numbers
	for patternLen := 1; patternLen <= maxDigits/2; patternLen++ {
		maxRepeats := maxDigits / patternLen
		if maxRepeats < 2 {
			continue
		}

		minPattern := 1
		if patternLen == 1 {
			minPattern = 0
		}
		maxPattern := 1
		for i := 0; i < patternLen; i++ {
			maxPattern *= 10
		}

		for pattern := minPattern; pattern < maxPattern; pattern++ {
			if countDigits(pattern) != patternLen && pattern != 0 {
				continue
			}

			for repeats := 2; repeats <= maxRepeats; repeats++ {
				num := generateRepeatedNumber(pattern, repeats)

				for _, r := range ranges {
					if num >= r.start && num <= r.end {
						validNumbers[num] = true
						break
					}
				}
			}
		}
	}

	total := 0
	for num := range validNumbers {
		total += num
	}

	return total
}

func hasRepeatingPattern(n int, buf []byte) bool {
	// Convert to string in-place
	s := strconv.AppendInt(buf[:0], int64(n), 10)
	sLen := len(s)

	// Try all possible pattern lengths from 1 to half the string length
	for patternLen := 1; patternLen <= sLen/2; patternLen++ {
		if sLen%patternLen != 0 {
			continue
		}

		repeats := sLen / patternLen
		matched := true

		// Compare bytes directly without creating substrings
		for j := 1; j < repeats; j++ {
			offset := j * patternLen
			for k := 0; k < patternLen; k++ {
				if s[k] != s[offset+k] {
					matched = false
					break
				}
			}
			if !matched {
				break
			}
		}

		if matched {
			return true
		}
	}

	return false
}

type Range struct {
	start int
	end   int
}

// Generate a number by repeating a pattern
func generateRepeatedNumber(pattern, repeats int) int {
	if repeats == 0 {
		return 0
	}

	// Calculate the multiplier for the pattern
	patternStr := strconv.Itoa(pattern)
	patternLen := len(patternStr)

	result := 0
	multiplier := 1

	// Build the number from right to left
	for i := 0; i < repeats; i++ {
		result += pattern * multiplier
		// Shift multiplier by the number of digits in pattern
		for j := 0; j < patternLen; j++ {
			multiplier *= 10
		}
	}

	return result
}

// Count digits in a number
func countDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n > 0 {
		count++
		n /= 10
	}
	return count
}
