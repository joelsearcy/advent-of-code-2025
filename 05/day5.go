package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// define a generic interval range type
type IntInterval struct {
	min, max int64
}

type IntIntervalTreeNode struct {
	value       IntInterval
	left, right *IntIntervalTreeNode
}

// Build a binary search interval tree constructed from a list of integer intervals
func buildBinaryIntervalTree(ranges []IntInterval, start, end int) *IntIntervalTreeNode {
	if start > end {
		return nil
	}

	mid := (start + end) / 2
	node := &IntIntervalTreeNode{
		value: ranges[mid],
	}
	node.left = buildBinaryIntervalTree(ranges, start, mid-1)
	node.right = buildBinaryIntervalTree(ranges, mid+1, end)
	return node
}

// Check if a value is contained within any of the intervals in the tree
func (node *IntIntervalTreeNode) contains(value int64) bool {
	if node == nil {
		return false
	} else if value >= node.value.min && value <= node.value.max {
		return true
	} else if node.left != nil && value < node.value.min {
		return node.left.contains(value)
	} else if node.right != nil && value > node.value.max {
		return node.right.contains(value)
	}
	return false
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	// split the input data into lines
	parts := strings.Split(strings.TrimSpace(string(data)), "\n\n")
	validRanges := strings.Split(parts[0], "\n")
	ingredients := strings.Split(parts[1], "\n")

	startTime := time.Now()
	total := 0

	// Parse valid integer ranges into an array of Interval[int] to build the interval tree
	// Ranges are separated by a dash '-'
	var ranges []IntInterval
	for _, line := range validRanges {
		var min, max int64
		fmt.Sscanf(line, "%d-%d", &min, &max)
		ranges = append(ranges, IntInterval{min: min, max: max})
	}

	// sort ranges by min value ascending, then max value descending
	sort.Slice(ranges, func(i, j int) bool {
		if ranges[i].min == ranges[j].min {
			return ranges[i].max > ranges[j].max
		}
		return ranges[i].min < ranges[j].min
	})

	// reduce/flatten overlapping ranges
	flattendRanges := []IntInterval{}
	for _, r := range ranges {
		if len(flattendRanges) == 0 {
			flattendRanges = append(flattendRanges, r)
			continue
		}
		last := &flattendRanges[len(flattendRanges)-1]
		if r.min <= last.max {
			if r.max > last.max {
				last.max = r.max
			}
		} else {
			flattendRanges = append(flattendRanges, r)
		}
	}

	// Build the interval tree from the flattened ranges
	intervalTree := buildBinaryIntervalTree(flattendRanges, 0, len(flattendRanges)-1)

	// Check each ingredient against the interval tree
	for _, line := range ingredients {
		var ingredient int64
		fmt.Sscanf(line, "%d", &ingredient)
		if intervalTree.contains(ingredient) {
			total++
		}
	}

	// part 2: sum range differences
	var totalValidIds int64
	for _, r := range flattendRanges {
		totalValidIds += r.max - r.min + 1
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Part 1 - valid ingredient IDs count: %d\n", total)
	fmt.Printf("Part 2 - Total potentially valid ingredient IDs: %d\n", totalValidIds)
	fmt.Printf("Execution time: %s\n", elapsed)
}
