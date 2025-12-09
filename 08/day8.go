package main

import (
	"cmp"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Point3D struct {
	X int
	Y int
	Z int
}

func (point Point3D) StraightlineDistance(other Point3D) float64 {
	dx := point.X - other.X
	dy := point.Y - other.Y
	dz := point.Z - other.Z
	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

type Edge struct {
	A        int
	B        int
	Distance float64
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	startTime := time.Now()
	// limit := 10
	// limit := 1000
	limit := math.MaxInt

	points := make([]Point3D, 0, len(lines))
	for _, line := range lines {
		coords := strings.Split(line, ",")
		if len(coords) != 3 {
			panic("invalid coordinate line: " + line)
		}
		x, err1 := strconv.Atoi(coords[0])
		y, err2 := strconv.Atoi(coords[1])
		z, err3 := strconv.Atoi(coords[2])
		if err1 != nil || err2 != nil || err3 != nil {
			panic("invalid coordinate values: " + line)
		}
		points = append(points, Point3D{X: x, Y: y, Z: z})
	}

	// Generate all edges
	edges := make([]Edge, 0)
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dist := points[i].StraightlineDistance(points[j])
			edges = append(edges, Edge{A: i, B: j, Distance: dist})
		}
	}

	// Sort edges by distance
	slices.SortFunc(edges, func(a, b Edge) int {
		return cmp.Compare(a.Distance, b.Distance)
	})

	// Build rooted "graphs" using Union-Find algorithm
	uf := NewUnionFind(len(points))
	count := 0
	lastEdge := Edge{}
	for i, edge := range edges {
		// fmt.Printf("Processing edge: %+v\n", edge)
		if i >= limit {
			break
		}
		if !uf.Connected(edge.A, edge.B) {
			// fmt.Printf("  Connecting points %d and %d\n", edge.A, edge.B)
			uf.Union(edge.A, edge.B)
			count++
			lastEdge = edge
		}
	}

	sizes := uf.TopNSizes(3)
	fmt.Printf("Sizes of (up to) three largest groups: %v\n", sizes)
	fmt.Printf("X's of last edge: %d, %d; product: %d\n", points[lastEdge.A].X, points[lastEdge.B].X, points[lastEdge.A].X*points[lastEdge.B].X)

	// Calculate product of three largest
	total := int64(1)
	if len(sizes) > 3 {
		sizes = sizes[:3]
	}
	for _, size := range sizes {
		total *= int64(size)
	}

	fmt.Printf("\nProduct of (up to) three largest groups: %d\n", total)
	fmt.Printf("Execution time: %s\n", time.Since(startTime))
}
