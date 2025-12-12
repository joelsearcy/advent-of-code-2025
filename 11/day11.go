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

	// line format: `sdf: you bla out`
	// before the colon is the machine/node name; afterwards are the machin/node names that output flows to
	// one of the machine names is `you`, representing the machines that you have direct access to
	// map paths from `you` to `out`.
	// `out` is a termination point.

	// considerations:
	//	1. Are there any circular paths?

	startTime := time.Now()

	// generate a graph of the input... map[string][]string?
	// parse lines into a map:
	graph := make(map[string][]string, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		node := parts[0]
		connections := strings.Split(parts[1], " ")
		graph[node] = connections
	}

	// Part 1: find every path from `you` to `out`, and output the count
	cache1 := make(map[string]int)
	part1Total := countPaths(graph, "you", "out", make(map[string]bool), true, true, cache1)

	// Part 2: find all paths from `svr` to `out` that also pass through `dac` and `fft` (in any order)
	cache2 := make(map[string]int)
	part2Total := countPaths(graph, "svr", "out", make(map[string]bool), false, false, cache2)

	fmt.Printf("Part 1 answer: %d\n", part1Total)
	fmt.Printf("Part 2 answer: %d\n", part2Total)
	fmt.Printf("Execution time: %s\n", time.Since(startTime))
}

func countPaths(graph map[string][]string, startNode string, endNode string, visited map[string]bool, foundDac bool, foundFft bool, cache map[string]int) int {
	if startNode == endNode {
		if foundDac && foundFft {
			return 1
		}
		return 0
	}

	// Update flags for required nodes
	if startNode == "dac" {
		foundDac = true
	}
	if startNode == "fft" {
		foundFft = true
	}

	// Check cache - key is node + state of required flags
	cacheKey := fmt.Sprintf("%s:%t:%t", startNode, foundDac, foundFft)
	if cachedCount, ok := cache[cacheKey]; ok {
		return cachedCount
	}

	linkedNodes, ok := graph[startNode]
	if !ok {
		return 0
	}

	visited[startNode] = true

	pathCount := 0
	for _, nextNode := range linkedNodes {
		if visited[nextNode] {
			continue
		}
		pathCount += countPaths(graph, nextNode, endNode, visited, foundDac, foundFft, cache)
	}

	delete(visited, startNode)

	cache[cacheKey] = pathCount
	return pathCount
}
