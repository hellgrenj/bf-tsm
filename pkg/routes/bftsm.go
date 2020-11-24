package routes

import (
	"math"
)

// Point represents a point in space.
type Point struct {
	Label string
	X     float64
	Y     float64
}

// NewPoint returns a Point based on X and Y positions on a graph.
func NewPoint(l string, x float64, y float64) Point {
	return Point{l, x, y}
}

// Distance finds the length of the hypotenuse between two points.
// Forumula is the square root of (x2 - x1)^2 + (y2 - y1)^2
func (p Point) Distance(p2 Point) float64 {
	first := math.Pow(float64(p2.X-p.X), 2)
	second := math.Pow(float64(p2.Y-p.Y), 2)
	return math.Sqrt(first + second)
}

// Heaps Algorithm: https://en.wikipedia.org/wiki/Heap%27s_algorithm
func permutations(arr []Point) [][]Point {
	var generate func([]Point, int)
	res := [][]Point{}

	generate = func(arr []Point, n int) {
		if n == 1 {
			tmp := make([]Point, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				generate(arr, n-1)
				if n%2 == 0 { // even
					arr[0], arr[n-1] = arr[n-1], arr[0] // swap first item in list with item at index n-1
				} else { // odd
					arr[i], arr[n-1] = arr[n-1], arr[i] // swap item at index i with item at index n-1
				}
			}
		}
	}
	generate(arr, len(arr))
	return res
}

// OptimalRoute contains the points, the score and the number of permutations compared
type OptimalRoute struct {
	Points           []Point
	Score            int
	NoOfPermutations int
}

// OptimalPath calculates and returns the optimal path
func OptimalPath(points []Point) OptimalRoute {
	permutations := permutations(points)
	var optimalRoute OptimalRoute
	for _, points := range permutations {

		pScore := 0
		for i := 0; i < len(points); i++ {

			if (i + 1) < len(points) {
				// we can also include other factors besides distance to calculate the score...
				pScore += int(points[i].Distance(points[i+1]))
			} else {
				// include the last leg.. going home..
				pScore += int(points[i].Distance(points[0]))
			}

		}
		if pScore < optimalRoute.Score || optimalRoute.Score == 0 {
			optimalRoute = OptimalRoute{Points: points, Score: pScore, NoOfPermutations: len(permutations)}
		}

	}
	return optimalRoute
}
