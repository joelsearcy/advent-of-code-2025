package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Area struct {
	Width         int
	Height        int
	PresentCounts map[int]int // shape index to count
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	startTime := time.Now()

	// parsing lines will have multiple phases:
	// first section are indexed lines with `<index>:` followed by lines that represent a grid of values (`.` for empty, `#` for part of the shape) terminated by a blank line
	// second section is a collection of lines defining an area to pack shapes into, starting with `<width>x<height>: ` then space-separated list of shape counts by shape indexes order to pack into the area

	// shapes can be rotated and flipped to fit into the area (up to 8 orientations each), without overlapping other shapes or going out of bounds
	// goal is to find if the shapes can all fit into the area as defined by the second section

	// split input into sections by blank lines, with the last section being the area definition
	sections := strings.Split(strings.TrimSpace(string(data)), "\n\n")
	rawAreas := strings.Split(sections[len(sections)-1], "\n")
	sections = sections[:len(sections)-1]
	presentShapes := make([][][]int, 0, len(sections))

	// parse shapes
	for _, present := range sections {
		lines := strings.Split(strings.TrimSpace(present), "\n")
		shape := make([][]int, len(lines)-1)
		for y, line := range lines[1:] {
			shape[y] = make([]int, len(line))
			for x, char := range line {
				if char == '#' {
					shape[y][x] = 1
				} else {
					shape[y][x] = 0
				}
			}
		}
		presentShapes = append(presentShapes, shape)
	}

	// parse area definitions
	areas := make([]Area, 0, len(rawAreas))
	for _, rawArea := range rawAreas {
		parts := strings.Split(rawArea, ": ")
		var width, height int
		fmt.Sscanf(parts[0], "%dx%d", &width, &height)
		counts := make(map[int]int)
		countParts := strings.Split(parts[1], " ")
		for i, countStr := range countParts {
			count, _ := strconv.Atoi(countStr)
			counts[i] = count
		}
		areas = append(areas, Area{
			Width:         width,
			Height:        height,
			PresentCounts: counts,
		})
	}

	// Part 1: For each area, determine if the shapes can fit into the area as defined
	// Process areas in parallel
	var wg sync.WaitGroup
	results := make([]bool, len(areas))
	for i, area := range areas {
		wg.Add(1)
		go func(idx int, a Area) {
			defer wg.Done()
			results[idx] = canPackShapes(a, presentShapes)

			//fmt.Printf("Area %d (%dx%d) processed: can pack = %v\n", idx, a.Width, a.Height, results[idx])
		}(i, area)
	}
	wg.Wait()

	part1Total := 0
	for _, result := range results {
		if result {
			part1Total++
		}
	}

	fmt.Printf("Part 1 Total: %d of %d\n", part1Total, len(areas))
	fmt.Printf("Execution time: %s\n", time.Since(startTime))
}

type Transform func(x, y, width, height int) (int, int)

// Identity (no transformation): (x,y) -> (x,y)
func identity(x, y, width, height int) (int, int) {
	return x, y
}

// Rotate 90° clockwise: (x,y) -> (y, width-1-x)
func rotate90CW(x, y, width, height int) (int, int) {
	return y, width - 1 - x
}

// Rotate 180°: (x,y) -> (width-1-x, height-1-y)
func rotate180(x, y, width, height int) (int, int) {
	return width - 1 - x, height - 1 - y
}

// Rotate 270° clockwise: (x,y) -> (height-1-y, x)
func rotate270CW(x, y, width, height int) (int, int) {
	return height - 1 - y, x
}

// Flip horizontal: (x,y) -> (width-1-x, y)
func flipH(x, y, width, height int) (int, int) {
	return width - 1 - x, y
}

// Flip vertical: (x,y) -> (x, height-1-y)
func flipV(x, y, width, height int) (int, int) {
	return x, height - 1 - y
}

// Get all unique orientations of a shape (deduplicated)
func getAllOrientations(shape [][]int) [][][]int {
	seen := make(map[string]bool)
	orientations := [][][]int{}

	transforms := []Transform{identity, rotate90CW, rotate180, rotate270CW}
	flips := []Transform{identity, flipH, flipV}

	for _, rot := range transforms {
		for _, flip := range flips {
			oriented := applyTransform(shape, func(x, y, w, h int) (int, int) {
				x1, y1 := rot(x, y, w, h)
				return flip(x1, y1, w, h)
			})
			key := shapeToString(oriented)
			if !seen[key] {
				seen[key] = true
				orientations = append(orientations, oriented)
			}
		}
	}

	return orientations
}

// Apply transformation to create a new shape matrix
func applyTransform(shape [][]int, transform Transform) [][]int {
	h, w := len(shape), len(shape[0])

	// Determine new dimensions after transform
	var newH, newW int
	x1, y1 := transform(0, 0, w, h)
	x2, y2 := transform(w-1, h-1, w, h)
	newW = abs(x2-x1) + 1
	newH = abs(y2-y1) + 1

	// Find min coordinates to normalize
	minX, minY := w, h
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			tx, ty := transform(x, y, w, h)
			if tx < minX {
				minX = tx
			}
			if ty < minY {
				minY = ty
			}
		}
	}

	result := make([][]int, newH)
	for i := range result {
		result[i] = make([]int, newW)
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if shape[y][x] == 1 {
				tx, ty := transform(x, y, w, h)
				result[ty-minY][tx-minX] = 1
			}
		}
	}

	return result
}

func shapeToString(shape [][]int) string {
	var sb strings.Builder
	for _, row := range shape {
		for _, val := range row {
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteString("|")
	}
	return sb.String()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Check if a shape can be placed at position (x, y) on the grid
func canPlaceShape(grid [][]int, shape [][]int, x, y int) bool {
	h, w := len(shape), len(shape[0])
	gridH, gridW := len(grid), len(grid[0])

	// Check bounds
	if y+h > gridH || x+w > gridW {
		return false
	}

	// Check overlap
	for sy := 0; sy < h; sy++ {
		for sx := 0; sx < w; sx++ {
			if shape[sy][sx] == 1 && grid[y+sy][x+sx] == 1 {
				return false
			}
		}
	}

	return true
}

// Place a shape on the grid (modifies grid in place)
func placeShape(grid [][]int, shape [][]int, x, y int) {
	h, w := len(shape), len(shape[0])
	for sy := 0; sy < h; sy++ {
		for sx := 0; sx < w; sx++ {
			if shape[sy][sx] == 1 {
				grid[y+sy][x+sx] = 1
			}
		}
	}
}

// Remove a shape from the grid (for backtracking)
func removeShape(grid [][]int, shape [][]int, x, y int) {
	h, w := len(shape), len(shape[0])
	for sy := 0; sy < h; sy++ {
		for sx := 0; sx < w; sx++ {
			if shape[sy][sx] == 1 {
				grid[y+sy][x+sx] = 0
			}
		}
	}
}

// Backtracking solver
func solveBacktrack(grid [][]int, remainingShapes []int, orientations [][][][]int) bool {
	// Base case: no more shapes to place
	if len(remainingShapes) == 0 {
		return true
	}

	// Take next shape to place
	shapeIdx := remainingShapes[0]
	remaining := remainingShapes[1:]

	// Early termination: check if remaining area is sufficient
	if len(remaining) > 0 {
		emptySpaces := countEmptySpaces(grid)
		remainingArea := 0
		for _, idx := range remaining {
			// Count cells in first orientation (all have same area)
			if len(orientations[idx]) > 0 {
				remainingArea += getShapeArea(orientations[idx][0])
			}
		}
		if remainingArea > emptySpaces {
			return false
		}
	}

	// Try all orientations
	for _, oriented := range orientations[shapeIdx] {
		h, w := len(oriented), len(oriented[0])
		gridH, gridW := len(grid), len(grid[0])

		// Try all positions
		for y := 0; y <= gridH-h; y++ {
			for x := 0; x <= gridW-w; x++ {
				if canPlaceShape(grid, oriented, x, y) {
					// Place and recurse
					placeShape(grid, oriented, x, y)

					if solveBacktrack(grid, remaining, orientations) {
						return true
					}

					// Backtrack
					removeShape(grid, oriented, x, y)
				}
			}
		}
	}

	return false
}

// Main packing function
func canPackShapes(area Area, presentShapes [][][]int) bool {
	// Pre-compute all orientations for each shape
	allOrientations := make([][][][]int, len(presentShapes))
	for i, shape := range presentShapes {
		allOrientations[i] = getAllOrientations(shape)
	}

	// Build list of shapes to place (with counts)
	type ShapeInfo struct {
		Idx  int
		Area int
	}
	shapeInfos := []ShapeInfo{}
	totalShapeArea := 0

	for shapeIdx, count := range area.PresentCounts {
		shapeArea := getShapeArea(presentShapes[shapeIdx])
		for i := 0; i < count; i++ {
			shapeInfos = append(shapeInfos, ShapeInfo{Idx: shapeIdx, Area: shapeArea})
			totalShapeArea += shapeArea
		}
	}

	// Early pruning: if total shape area > grid area, impossible
	gridArea := area.Width * area.Height
	if totalShapeArea > gridArea {
		return false
	}

	// Sort shapes by area (largest first - most constrained)
	sort.Slice(shapeInfos, func(i, j int) bool {
		return shapeInfos[i].Area > shapeInfos[j].Area
	})

	shapesToPlace := make([]int, len(shapeInfos))
	for i, info := range shapeInfos {
		shapesToPlace[i] = info.Idx
	}

	// Create empty grid
	grid := make([][]int, area.Height)
	for i := range grid {
		grid[i] = make([]int, area.Width)
	}

	// Solve with backtracking
	return solveBacktrack(grid, shapesToPlace, allOrientations)
}

// Calculate area of a shape (number of filled cells)
func getShapeArea(shape [][]int) int {
	area := 0
	for _, row := range shape {
		for _, cell := range row {
			if cell == 1 {
				area++
			}
		}
	}
	return area
}

// Count empty spaces in grid
func countEmptySpaces(grid [][]int) int {
	count := 0
	for _, row := range grid {
		for _, cell := range row {
			if cell == 0 {
				count++
			}
		}
	}
	return count
}
