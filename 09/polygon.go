package main

// assume that the polygon is simple (no self-intersections) and right-handed (points ordered clockwise)
type Polygon struct {
	Points [][2]int
}

func (poly *Polygon) Contains(point [2]int) bool {
	// Check if point is on the perimeter first
	if poly.IsOnPerimeter(point) {
		return true
	}

	// Ray casting algorithm: cast a ray from the point to infinity
	// and count how many polygon edges it crosses.
	// If the count is odd, the point is inside; if even, it's outside.
	n := len(poly.Points)
	inside := false

	for i := 0; i < n; i++ {
		p1 := poly.Points[i]
		p2 := poly.Points[(i+1)%n]

		// Use half-open interval [min, max) to handle vertices consistently
		// This ensures a vertex is only counted once when ray passes through it
		if (p1[1] <= point[1] && point[1] < p2[1]) || (p2[1] <= point[1] && point[1] < p1[1]) {
			// The edge crosses the horizontal ray from the point
			// Calculate the x-coordinate of the intersection
			xIntersection := float64(p1[0]) + float64(point[1]-p1[1])*float64(p2[0]-p1[0])/float64(p2[1]-p1[1])

			if float64(point[0]) < xIntersection {
				inside = !inside
			}
		}
	}
	return inside
}

func (poly *Polygon) IsOnPerimeter(point [2]int) bool {
	// Check if point lies on any edge of the polygon
	n := len(poly.Points)
	for i := 0; i < n; i++ {
		p1 := poly.Points[i]
		p2 := poly.Points[(i+1)%n]
		if pointOnSegment(p1, p2, point) {
			return true
		}
	}
	return false
}

func pointOnSegment(p1, p2, p [2]int) bool {
	// Check if p is collinear with p1-p2 and lies between them
	if (p[0]-p1[0])*(p2[1]-p1[1]) != (p[1]-p1[1])*(p2[0]-p1[0]) {
		return false // not collinear
	}
	// Check if p is within bounding box of p1-p2
	return p[0] >= min(p1[0], p2[0]) && p[0] <= max(p1[0], p2[0]) &&
		p[1] >= min(p1[1], p2[1]) && p[1] <= max(p1[1], p2[1])
}
