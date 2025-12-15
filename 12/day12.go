package main

import (
	"flag"
	"fmt"
	"math/bits"
	"os"
	"runtime/pprof"
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

// BitGrid represents a grid using uint64 bitsets for ultra-fast operations
type BitGrid struct {
	rows   []uint64
	width  int
	height int
}

// BitShape represents a shape using uint64 bitsets
type BitShape struct {
	rows   []uint64
	width  int
	height int
	area   int
}

// Convert shape to bitset
func shapeToBitset(shape [][]int) BitShape {
	h, w := len(shape), len(shape[0])
	rows := make([]uint64, h)
	area := 0

	for y := 0; y < h; y++ {
		var row uint64
		for x := 0; x < w; x++ {
			if shape[y][x] == 1 {
				row |= 1 << x
				area++
			}
		}
		rows[y] = row
	}

	return BitShape{rows: rows, width: w, height: h, area: area}
}

// Create empty bit grid
func newBitGrid(width, height int) BitGrid {
	return BitGrid{
		rows:   make([]uint64, height),
		width:  width,
		height: height,
	}
}

// Check if shape can be placed at position (x, y) - ULTRA FAST
func canPlaceBitShape(grid *BitGrid, shape *BitShape, x, y int) bool {
	// Bounds check
	if y+shape.height > grid.height || x+shape.width > grid.width {
		return false
	}

	// Check overlap using bitwise AND with early exit
	for sy := 0; sy < shape.height; sy++ {
		shapeRow := shape.rows[sy] << x
		if (grid.rows[y+sy] & shapeRow) != 0 {
			return false
		}
	}

	return true
}

// Place shape on grid - ULTRA FAST
func placeBitShape(grid *BitGrid, shape *BitShape, x, y int) {
	for sy := 0; sy < shape.height; sy++ {
		grid.rows[y+sy] |= shape.rows[sy] << x
	}
}

// Remove shape from grid - ULTRA FAST
func removeBitShape(grid *BitGrid, shape *BitShape, x, y int) {
	for sy := 0; sy < shape.height; sy++ {
		grid.rows[y+sy] ^= shape.rows[sy] << x
	}
}

// Count filled cells in grid using fast popcount
func countFilledCells(grid *BitGrid) int {
	count := 0
	for _, row := range grid.rows {
		count += bits.OnesCount64(row)
	}
	return count
}

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flag.String("memprofile", "", "write memory profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

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

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := pprof.WriteHeapProfile(f); err != nil {
			panic(err)
		}
	}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Backtracking solver with index-based recursion (zero allocations)
func solveBacktrackBitMRV(grid *BitGrid, shapesToPlace []int, startIdx int, orientations [][]BitShape, remainingArea int, emptySpaces int) bool {
	// Base case: placed all shapes
	if startIdx >= len(shapesToPlace) {
		return true
	}

	// Early termination: if remaining shapes can't fit in empty space
	if remainingArea > emptySpaces {
		return false
	}

	// Additional pruning: use popcount to verify actual empty space
	// Check every 3 levels to catch fragmentation early
	if startIdx%3 == 0 && startIdx > 0 {
		actualFilled := countFilledCells(grid)
		actualEmpty := grid.width*grid.height - actualFilled
		if remainingArea > actualEmpty {
			return false
		}
	}

	// Symmetry breaking: if grid is empty, only try top-left corner for first shape
	if startIdx == 0 {
		shapeIdx := shapesToPlace[0]

		for i := range orientations[shapeIdx] {
			oriented := &orientations[shapeIdx][i]
			if oriented.height <= grid.height && oriented.width <= grid.width {
				if canPlaceBitShape(grid, oriented, 0, 0) {
					placeBitShape(grid, oriented, 0, 0)
					if solveBacktrackBitMRV(grid, shapesToPlace, 1, orientations, remainingArea-oriented.area, emptySpaces-oriented.area) {
						return true
					}
					removeBitShape(grid, oriented, 0, 0)
				}
			}
		}
		return false
	}

	// Get current shape to place
	shapeIdx := shapesToPlace[startIdx]

	// Try all orientations and positions
	for i := range orientations[shapeIdx] {
		oriented := &orientations[shapeIdx][i]

		if oriented.height > grid.height || oriented.width > grid.width {
			continue
		}

		maxY := grid.height - oriented.height
		maxX := grid.width - oriented.width

		for y := 0; y <= maxY; y++ {
			row := grid.rows[y]

			// Skip completely filled rows
			if row == (uint64(1)<<grid.width)-1 {
				continue
			}

			for x := 0; x <= maxX; x++ {
				// Use TrailingZeros64 to skip filled positions
				if x < 64 && (row&(1<<x)) != 0 {
					// Position is filled - find next empty position
					remaining := ^row >> (x + 1) // Invert and shift past current position
					if remaining == 0 {
						break // No more empty cells
					}
					skip := bits.TrailingZeros64(remaining)
					x += skip + 1
					if x > maxX {
						break
					}
				}

				if canPlaceBitShape(grid, oriented, x, y) {
					placeBitShape(grid, oriented, x, y)

					if solveBacktrackBitMRV(grid, shapesToPlace, startIdx+1, orientations, remainingArea-oriented.area, emptySpaces-oriented.area) {
						return true
					}

					removeBitShape(grid, oriented, x, y)
				}
			}
		}
	}

	return false
}

// Cache for orientations to avoid recomputing across areas
var orientationCache = struct {
	sync.RWMutex
	cache map[string][]BitShape
}{cache: make(map[string][]BitShape)}

// Get all orientations as bitsets
func getAllOrientationsBit(shape [][]int) []BitShape {
	orientations := getAllOrientations(shape)
	bitOrientations := make([]BitShape, len(orientations))

	for i, oriented := range orientations {
		bitOrientations[i] = shapeToBitset(oriented)
	}

	return bitOrientations
}

// Main packing function with bitset optimization
func canPackShapes(area Area, presentShapes [][][]int) bool {
	// Pre-compute all orientations as bitsets (with caching)
	allOrientations := make([][]BitShape, len(presentShapes))
	for i, shape := range presentShapes {
		key := shapeToString(shape)

		orientationCache.RLock()
		cached, ok := orientationCache.cache[key]
		orientationCache.RUnlock()

		if ok {
			allOrientations[i] = cached
		} else {
			orientations := getAllOrientationsBit(shape)
			orientationCache.Lock()
			orientationCache.cache[key] = orientations
			orientationCache.Unlock()
			allOrientations[i] = orientations
		}
	}

	// Build list of shapes to place (with counts)
	type ShapeInfo struct {
		Idx  int
		Area int
	}

	// Pre-allocate to avoid reallocation
	totalCount := 0
	for _, count := range area.PresentCounts {
		totalCount += count
	}
	shapeInfos := make([]ShapeInfo, 0, totalCount)
	totalShapeArea := 0

	for shapeIdx, count := range area.PresentCounts {
		shapeArea := getShapeArea(presentShapes[shapeIdx])
		for i := 0; i < count; i++ {
			shapeInfos = append(shapeInfos, ShapeInfo{
				Idx:  shapeIdx,
				Area: shapeArea,
			})
			totalShapeArea += shapeArea
		}
	}

	// Early pruning: if total shape area > grid area, impossible
	gridArea := area.Width * area.Height
	if totalShapeArea > gridArea {
		return false
	}

	// Sort shapes smallest first (proven fastest strategy)
	sort.Slice(shapeInfos, func(i, j int) bool {
		return shapeInfos[i].Area < shapeInfos[j].Area
	})

	// Pre-allocated with capacity
	shapesToPlace := make([]int, len(shapeInfos))
	remainingSizes := totalShapeArea
	for i, info := range shapeInfos {
		shapesToPlace[i] = info.Idx
	}

	// Create empty bitgrid
	grid := newBitGrid(area.Width, area.Height)

	// Solve with bitset backtracking (index-based, zero allocations)
	return solveBacktrackBitMRV(&grid, shapesToPlace, 0, allOrientations, remainingSizes, gridArea)
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
