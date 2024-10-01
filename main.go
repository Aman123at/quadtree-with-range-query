package main

import (
	"fmt"
)

// Point represents a 2D point
type Point struct {
	X, Y float64
}

// Rect represents a rectangle
type Rect struct {
	X, Y, Width, Height float64
}

// QuadTree represents the quadtree data structure
type QuadTree struct {
	Boundary     Rect
	Capacity     int
	Points       []Point
	NorthWest    *QuadTree
	NorthEast    *QuadTree
	SouthWest    *QuadTree
	SouthEast    *QuadTree
	IsSubdivided bool
}

// NewQuadTree creates a new QuadTree
func NewQuadTree(boundary Rect, capacity int) *QuadTree {
	return &QuadTree{
		Boundary: boundary,
		Capacity: capacity,
		Points:   make([]Point, 0),
	}
}

// Insert adds a point to the QuadTree
func (qt *QuadTree) Insert(p Point) bool {
	if !qt.Boundary.Contains(p) {
		return false
	}

	if len(qt.Points) < qt.Capacity && !qt.IsSubdivided {
		qt.Points = append(qt.Points, p)
		return true
	}

	if !qt.IsSubdivided {
		qt.Subdivide()
	}

	if qt.NorthWest.Insert(p) {
		return true
	}
	if qt.NorthEast.Insert(p) {
		return true
	}
	if qt.SouthWest.Insert(p) {
		return true
	}
	if qt.SouthEast.Insert(p) {
		return true
	}

	return false
}

// Subdivide splits the QuadTree into four quadrants
func (qt *QuadTree) Subdivide() {
	x := qt.Boundary.X
	y := qt.Boundary.Y
	w := qt.Boundary.Width / 2
	h := qt.Boundary.Height / 2

	qt.NorthWest = NewQuadTree(Rect{x, y, w, h}, qt.Capacity)
	qt.NorthEast = NewQuadTree(Rect{x + w, y, w, h}, qt.Capacity)
	qt.SouthWest = NewQuadTree(Rect{x, y + h, w, h}, qt.Capacity)
	qt.SouthEast = NewQuadTree(Rect{x + w, y + h, w, h}, qt.Capacity)

	qt.IsSubdivided = true

	for _, p := range qt.Points {
		qt.Insert(p)
	}
	qt.Points = nil
}

// Search finds a point in the QuadTree
func (qt *QuadTree) Search(p Point) bool {
	if !qt.Boundary.Contains(p) {
		return false
	}

	for _, point := range qt.Points {
		if point == p {
			return true
		}
	}

	if qt.IsSubdivided {
		return qt.NorthWest.Search(p) ||
			qt.NorthEast.Search(p) ||
			qt.SouthWest.Search(p) ||
			qt.SouthEast.Search(p)
	}

	return false
}

// Query returns all points within a given rectangle
func (qt *QuadTree) Query(range_ Rect) []Point {
	found := make([]Point, 0)

	if !qt.Boundary.Intersects(range_) {
		return found
	}

	for _, p := range qt.Points {
		if range_.Contains(p) {
			found = append(found, p)
		}
	}

	if qt.IsSubdivided {
		found = append(found, qt.NorthWest.Query(range_)...)
		found = append(found, qt.NorthEast.Query(range_)...)
		found = append(found, qt.SouthWest.Query(range_)...)
		found = append(found, qt.SouthEast.Query(range_)...)
	}

	return found
}

// Delete removes a point from the QuadTree
func (qt *QuadTree) Delete(p Point) bool {
	if !qt.Boundary.Contains(p) {
		return false
	}

	for i, point := range qt.Points {
		if point == p {
			qt.Points = append(qt.Points[:i], qt.Points[i+1:]...)
			return true
		}
	}

	if qt.IsSubdivided {
		if qt.NorthWest.Delete(p) {
			return true
		}
		if qt.NorthEast.Delete(p) {
			return true
		}
		if qt.SouthWest.Delete(p) {
			return true
		}
		if qt.SouthEast.Delete(p) {
			return true
		}
	}

	return false
}

// Contains checks if a point is within the rectangle
func (r Rect) Contains(p Point) bool {
	return p.X >= r.X && p.X < r.X+r.Width &&
		p.Y >= r.Y && p.Y < r.Y+r.Height
}

// Intersects checks if two rectangles intersect
func (r Rect) Intersects(other Rect) bool {
	return !(other.X > r.X+r.Width ||
		other.X+other.Width < r.X ||
		other.Y > r.Y+r.Height ||
		other.Y+other.Height < r.Y)
}

// MapPoint represents a point on the map with additional data
type MapPoint struct {
	Point
	Data interface{}
}

// Map represents the zoomable map
type Map struct {
	QuadTree       *QuadTree
	ZoomLevel      int
	CenterX        float64
	CenterY        float64
	ViewportWidth  float64
	ViewportHeight float64
}

// NewMap creates a new Map
func NewMap(boundary Rect, capacity int) *Map {
	return &Map{
		QuadTree:       NewQuadTree(boundary, capacity),
		ZoomLevel:      0,
		CenterX:        boundary.X + boundary.Width/2,
		CenterY:        boundary.Y + boundary.Height/2,
		ViewportWidth:  boundary.Width,
		ViewportHeight: boundary.Height,
	}
}

// AddPoint adds a point to the map
func (m *Map) AddPoint(p MapPoint) {
	m.QuadTree.Insert(p.Point)
}

// ZoomIn increases the zoom level and adjusts the viewport
func (m *Map) ZoomIn() {
	m.ZoomLevel++
	m.ViewportWidth /= 2
	m.ViewportHeight /= 2
}

// ZoomOut decreases the zoom level and adjusts the viewport
func (m *Map) ZoomOut() {
	if m.ZoomLevel > 0 {
		m.ZoomLevel--
		m.ViewportWidth *= 2
		m.ViewportHeight *= 2
	}
}

// Pan moves the center of the viewport
func (m *Map) Pan(dx, dy float64) {
	m.CenterX += dx
	m.CenterY += dy
}

// GetVisiblePoints returns all points visible in the current viewport
func (m *Map) GetVisiblePoints() []MapPoint {
	viewportRect := Rect{
		X:      m.CenterX - m.ViewportWidth/2,
		Y:      m.CenterY - m.ViewportHeight/2,
		Width:  m.ViewportWidth,
		Height: m.ViewportHeight,
	}
	points := m.QuadTree.Query(viewportRect)

	mapPoints := make([]MapPoint, len(points))
	for i, p := range points {
		mapPoints[i] = MapPoint{Point: p}
	}
	return mapPoints
}

func main() {
	boundary := Rect{0, 0, 1000, 1000}
	map_new := NewMap(boundary, 4)

	// Add some points
	points := []MapPoint{
		{Point: Point{100, 100}, Data: "Point A"},
		{Point: Point{200, 200}, Data: "Point B"},
		{Point: Point{300, 300}, Data: "Point C"},
		{Point: Point{400, 400}, Data: "Point D"},
		{Point: Point{500, 500}, Data: "Point E"},
	}

	for _, p := range points {
		map_new.AddPoint(p)
	}

	// Initial view
	fmt.Println("Initial view:")
	fmt.Printf("Zoom level: %d\n", map_new.ZoomLevel)
	fmt.Printf("Visible points: %v\n", map_new.GetVisiblePoints())

	// Zoom in twice
	map_new.ZoomIn()
	map_new.ZoomIn()
	fmt.Println("\nAfter zooming in twice:")
	fmt.Printf("Zoom level: %d\n", map_new.ZoomLevel)
	fmt.Printf("Visible points: %v\n", map_new.GetVisiblePoints())

	// Pan the map
	map_new.Pan(100, 100)
	fmt.Println("\nAfter panning:")
	fmt.Printf("Center: (%f, %f)\n", map_new.CenterX, map_new.CenterY)
	fmt.Printf("Visible points: %v\n", map_new.GetVisiblePoints())

	// Zoom out once
	map_new.ZoomOut()
	fmt.Println("\nAfter zooming out once:")
	fmt.Printf("Zoom level: %d\n", map_new.ZoomLevel)
	fmt.Printf("Visible points: %v\n", map_new.GetVisiblePoints())
}
