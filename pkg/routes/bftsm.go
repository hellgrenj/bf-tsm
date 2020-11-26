package routes

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

// Point represents a point in space.
type Point struct {
	Label string
	X     float64
	Y     float64
}

// NewPoint returns a Point with a label based on X and Y positions on a graph.
func NewPoint(l string, x float64, y float64) Point {
	return Point{l, x, y}
}

// finds the length of the hypotenuse between two points.
// Forumula is the square root of (x2 - x1)^2 + (y2 - y1)^2
func (p Point) distance(p2 Point) float64 {
	first := math.Pow(float64(p2.X-p.X), 2)
	second := math.Pow(float64(p2.Y-p.Y), 2)
	return math.Sqrt(first + second)
}

// Heaps Algorithm: https://en.wikipedia.org/wiki/Heap%27s_algorithm
func permutations(arr []Point) [][]Point {
	permExecStart := time.Now()
	defer func() {
		permExecTime := time.Since(permExecStart)
		fmt.Printf("\ncalculating the permutations took %v ms\n", permExecTime.Milliseconds())
	}()
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

// OptimalRoute contains the points, the cost and the number of permutations compared
type OptimalRoute struct {
	Points               []Point
	Cost                 int
	NumberOfPermutations int
}

// OptimalPath calculates and returns the optimal path
func OptimalPath(points []Point) OptimalRoute {
	permutations := permutations(points)

	fmt.Printf("\nNumber of permutations %v\n", len(permutations))
	// calculate optimal route, in parallel if more then 300 000 permutations..
	if len(permutations) > 300000 {
		parallelOptimalRouteCalcStart := time.Now()
		numberOfCores := runtime.NumCPU()
		chunkSize := int(len(permutations) / numberOfCores)
		fmt.Printf("\ndistributing the work over the %v CPU cores\n", numberOfCores)
		var chans = make([]chan OptimalRoute, numberOfCores)
		for i := 0; i < numberOfCores; i++ {
			chans[i] = make(chan OptimalRoute)
			if i == 0 {
				// first chunk (take everything from 0 to chunkSize)
				go parallelOptimalRoute(permutations[i:chunkSize], chans[i])
			} else if i == (numberOfCores - 1) {
				// last chunk (take all that is left...)
				go parallelOptimalRoute(permutations[(chunkSize*i):], chans[i])
			} else {
				// chunks inbetween start and end
				go parallelOptimalRoute(permutations[(chunkSize*(i)):(chunkSize*(i+1))], chans[i])
			}
		}
		fmt.Println("now waiting for all channels to return their batch result")
		var optimalRoutes []OptimalRoute
		for _, v := range chans {
			or := <-v
			optimalRoutes = append(optimalRoutes, or)
		}
		var finalPerms [][]Point
		for _, o := range optimalRoutes {
			finalPerms = append(finalPerms, o.Points)
		}
		fmt.Println("now running the candidates")
		optimalRoute := calculateOptimalRoute(finalPerms)
		fmt.Printf("\nreturning the winner %v\n", optimalRoute)
		optimalRoute.NumberOfPermutations = len(permutations)
		parallelOptimalRouteCalcTime := time.Since(parallelOptimalRouteCalcStart)
		fmt.Printf("\n(parallel) optimal route calc took %v ms\n", parallelOptimalRouteCalcTime.Milliseconds())
		return optimalRoute
	}
	// if fewer permutations run everything in this (one) go routine...
	optimalRouteCalcStart := time.Now()
	optimalRoute := calculateOptimalRoute(permutations)
	optimalRouteCalcTime := time.Since(optimalRouteCalcStart)
	fmt.Printf("\noptimal route calc took %v ms\n", optimalRouteCalcTime.Milliseconds())
	return optimalRoute
}
func parallelOptimalRoute(permutations [][]Point, ch chan OptimalRoute) {
	optimalRoute := calculateOptimalRoute(permutations)
	ch <- optimalRoute
}
func calculateOptimalRoute(permutations [][]Point) OptimalRoute {
	var optimalRoute OptimalRoute
	for _, points := range permutations {
		cost := 0
		for i := range points {
			// we can also include other factors besides distance to calculate the cost...
			if (i + 1) < len(points) {
				cost += int(points[i].distance(points[i+1]))
			} else {
				// include the last leg.. going home..
				cost += int(points[i].distance(points[0]))
			}
		}
		if cost < optimalRoute.Cost || optimalRoute.Cost == 0 {
			optimalRoute = OptimalRoute{Points: points, Cost: cost, NumberOfPermutations: len(permutations)}
		}
	}
	return optimalRoute
}
