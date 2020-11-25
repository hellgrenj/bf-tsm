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
	permExecStart := time.Now()
	permutations := permutations(points)
	permExecTime := time.Since(permExecStart)
	fmt.Printf("\npermutations took %v ms\n", permExecTime.Milliseconds())

	// calculate optimal route, in parallel if more then 3 000 000 permutations..
	fmt.Printf("\nNumber of permutations %v\n", len(permutations))
	if len(permutations) > 300000 {
		numberOfCores := runtime.NumCPU()
		fmt.Printf("\nnumber of cores %v\n", numberOfCores)
		chunkSize := int(len(permutations) / numberOfCores)
		if numberOfCores > 100 {
			numberOfCores = 100
		}
		var chans [100]chan OptimalRoute // set to big engough... (now only supports 100 cores.. :) )
		for i := 0; i < numberOfCores; i++ {
			chans[i] = make(chan OptimalRoute)
			if i == 0 {
				fmt.Printf("\nfirst chunk take everything from 0 to chunkSize (i is %v, and chunkSize is %v\n", i, chunkSize)
				// first chunk (take everything from 0 to chunkSize)
				go parallelOptimalRoute(permutations[i:chunkSize], chans[i])
			} else if i == (numberOfCores - 1) {
				fmt.Printf("\nlast chunk (take all that is left...) (from index (chunkSize*i) which is %v\n", (chunkSize * i))
				// last chunk (take all that is left...)
				go parallelOptimalRoute(permutations[(chunkSize*i):], chans[i])
			} else {
				fmt.Printf("\nchunks inbetween start and end (chunkSize*(i) is %v, and chunkSize*(i+1) is %v\n", chunkSize*(i), chunkSize*(i+1))
				// chunks inbetween start and end
				go parallelOptimalRoute(permutations[(chunkSize*(i)):(chunkSize*(i+1))], chans[i])
			}
		}
		fmt.Println("now waiting for all channels to return their batch result")
		var optimalRoutes []OptimalRoute
		for _, v := range chans[:numberOfCores] { // disregard unused nil channels...
			or := <-v
			optimalRoutes = append(optimalRoutes, or)
		}
		var finalPerms [][]Point
		for _, o := range optimalRoutes {
			finalPerms = append(finalPerms, o.Points)
		}
		fmt.Println("now running the candidates")
		optimalRoute := optimalRoute(finalPerms)
		fmt.Printf("\nreturning the winner %v\n", optimalRoute)
		optimalRoute.NoOfPermutations = len(permutations)
		return optimalRoute
	}
	// if fewer permutations run everything in this (one) go routine...
	optimalRouteCalcStart := time.Now()
	optimalRoute := optimalRoute(permutations)
	optimalRouteCalcTime := time.Since(optimalRouteCalcStart)
	fmt.Printf("\noptimal route calc took %v ms\n", optimalRouteCalcTime.Milliseconds())
	return optimalRoute
}
func parallelOptimalRoute(permutations [][]Point, ch chan OptimalRoute) {
	or := optimalRoute(permutations)
	ch <- or
}
func optimalRoute(permutations [][]Point) OptimalRoute {
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
