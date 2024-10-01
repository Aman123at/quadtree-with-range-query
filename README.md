# QuadTree With Range Query

## Overview

This Go package demonstrates the implementation of a **QuadTree**, a tree data structure used for spatial indexing, and a Map, which represents a zoomable map using the QuadTree for efficiently managing and querying spatial points.

## QuadTree
A **QuadTree** is a data structure that recursively divides a 2D space into four quadrants (northwest, northeast, southwest, southeast). It is useful for efficient querying of points in 2D space and spatial subdivision.

## Map
The Map uses a QuadTree to store and manage spatial data (points on a map), allowing zooming and panning functionality for easy manipulation of the map view.

## QuadTree Structure

### Types

`Point`

Represents a 2D point with `X` and `Y` coordinates.
```go
type Point struct {
    X, Y float64
}
```

`Rect`

Represents a rectangle, used as boundaries for quadrants in the QuadTree.
```go
type Rect struct {
    X, Y, Width, Height float64
}
```

`QuadTree`

Represents the QuadTree structure. It stores points and subdivides the space when the number of points exceeds its capacity.
```go
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
```


## Functions

`NewQuadTree(boundary Rect, capacity int) *QuadTree`

Creates a new QuadTree with a specified boundary and capacity (the maximum number of points before subdivision).
```go
qt := NewQuadTree(Rect{0, 0, 100, 100}, 4)
```

`Insert(p Point) bool`

Inserts a point into the QuadTree. If the QuadTree is full, it subdivides into four quadrants.
```go
qt.Insert(Point{X: 10, Y: 20})
```


`Subdivide()`

Subdivides the QuadTree into four smaller quadrants: northwest, northeast, southwest, and southeast.
```go
qt.Subdivide()
```


`Search(p Point) bool`

Searches for a point in the QuadTree.
```go
exists := qt.Search(Point{X: 10, Y: 20})
```


`Query(range_ Rect) []Point`

Returns all points within a given rectangle.
```go
points := qt.Query(Rect{X: 0, Y: 0, Width: 50, Height: 50})
```


`Delete(p Point) bool`

Removes a point from the QuadTree.
```go
qt.Delete(Point{X: 10, Y: 20})
```


## Helper Functions

`Contains(p Point) bool`

Checks if a point lies within a rectangle.
```go
r.Contains(Point{X: 10, Y: 20})
```

`Intersects(other Rect) bool`

Checks if two rectangles intersect.
```go
r.Intersects(Rect{X: 50, Y: 50, Width: 100, Height: 100})
```


## Map Structure

### Types

`MapPoint`

Represents a point on the map with additional data.
```go
type MapPoint struct {
    Point
    Data interface{}
}
```

`Map`

Represents a map with zooming and panning functionality.
```go
type Map struct {
    QuadTree       *QuadTree
    ZoomLevel      int
    CenterX        float64
    CenterY        float64
    ViewportWidth  float64
    ViewportHeight float64
}
```

## Functions

`NewMap(boundary Rect, capacity int) *Map`

Creates a new map with the given boundary and QuadTree capacity.
```go
m := NewMap(Rect{0, 0, 1000, 1000}, 4)
```

`AddPoint(p MapPoint)`

Adds a point to the map.
```go
m.AddPoint(MapPoint{Point: Point{100, 100}, Data: "Point A"})
```

`ZoomIn()`

Zooms in by reducing the viewport size.
```go
m.ZoomIn()
```

`ZoomOut()`

Zooms out by increasing the viewport size.
```go
m.ZoomOut()
```

`Pan(dx, dy float64)`

Pans the map by moving the center.

```go
m.Pan(100, 100)
```

`GetVisiblePoints() []MapPoint`

Returns all points visible within the current viewport.

```go
visiblePoints := m.GetVisiblePoints()
```



















