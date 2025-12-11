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

	startTime := time.Now()
	maxArea := int64(0)

	points := make([][2]int, 0, len(lines))
	for _, line := range lines {
		coords := strings.Split(line, ",")
		if len(coords) != 2 {
			panic("invalid coordinate line: " + line)
		}
		x, err1 := strconv.Atoi(coords[0])
		y, err2 := strconv.Atoi(coords[1])
		if err1 != nil || err2 != nil {
			panic("invalid coordinate values: " + line)
		}
		points = append(points, [2]int{x, y})
	}

	// Part 1: Find the largest area of a rectangle defined by any two points
	for i := 0; i < len(points); i++ {
		area := int64(0)
		for j := i + 1; j < len(points); j++ {
			dx := int64(abs(points[j][0]-points[i][0])) + 1
			dy := int64(abs(points[j][1]-points[i][1])) + 1
			area = dx * dy
			if area > maxArea {
				maxArea = area
			}
		}
	}
	fmt.Printf("Part 1 - max area of a rectangle (ignoring perimeter): %d\n", maxArea)

	// sortedPoints := make([][2]int, len(points))
	// copy(sortedPoints, points)
	// // sort points by y, then by x
	// slices.SortFunc(sortedPoints, func(a, b [2]int) int {
	// 	if a[1] != b[1] {
	// 		return a[1] - b[1]
	// 	}
	// 	return a[0] - b[0]
	// })
	polygon := Polygon{Points: points}

	// Part 2: Find the largest area of a rectangle contained within perimeter points
	maxArea = 0
	for i := 0; i < len(points); i++ {
		area := int64(0)
		for j := i + 1; j < len(points); j++ {
			minX := min(points[i][0], points[j][0])
			maxX := max(points[i][0], points[j][0])
			minY := min(points[i][1], points[j][1])
			maxY := max(points[i][1], points[j][1])

			contained := true
			// check that no polygon edges cross into or through the rectangle
			// if the polygon edges where diagnal, this would be more complex
			for k := 0; k < len(points); k++ {
				// points p->q are the current polygon edge being considered
				p := points[k]
				q := points[(k+1)%len(points)]

				// check horizontal edges intersecting the rectangle
				if p[1] == q[1] && p[1] > minY && p[1] < maxY && min(p[0], q[0]) < maxX && max(p[0], q[0]) > minX {
					contained = false
					break
				}
				// check vertical edges intersecting the rectangle
				if p[0] == q[0] && p[0] > minX && p[0] < maxX && min(p[1], q[1]) < maxY && max(p[1], q[1]) > minY {
					contained = false
					break
				}
			}
			if !contained {
				continue
			}
			// check that all four corners are within the polygon
			rectPoints := [][2]int{{minX, minY}, {maxX, minY}, {maxX, maxY}, {minX, maxY}}
			for _, rp := range rectPoints {
				if !polygon.Contains(rp) {
					contained = false
					break
				}
			}
			if !contained {
				continue
			}

			dx := maxX - minX
			dy := maxY - minY
			area = int64((dx + 1) * (dy + 1))
			if area > maxArea {
				maxArea = area
			}
		}
	}
	fmt.Printf("Part 2 - max area of a rectangle (on/within perimeter): %d\n", maxArea)

	fmt.Printf("Execution time: %s\n", time.Since(startTime))
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
