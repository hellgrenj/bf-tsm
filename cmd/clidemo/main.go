package main

import (
	"fmt"
	"time"

	"github.com/hellgrenj/bf-tsm/pkg/routes"
)

func prettyPrintRoute(route []routes.Point) {
	for i, p := range route {
		if i == len(route)-1 {
			fmt.Printf("%v => back to %v\n", p.Label, route[0].Label)
		} else {
			fmt.Printf("%v, ", p.Label)
		}

	}
}
func main() {
	amsterdam := routes.NewPoint("Amsterdam", 52.377956, 4.897070)
	berlin := routes.NewPoint("Berlin", 52.520008, 13.404954)
	kiruna := routes.NewPoint("Kiruna", 67.85000, 20.23000)
	goteborg := routes.NewPoint("Göteborg", 57.69000, 11.89000)
	johannesburg := routes.NewPoint("Johannes Burg", -26.195246, 28.034088)
	newyork := routes.NewPoint("New york", 40.730610, -73.935242)
	havana := routes.NewPoint("Havana", 23.113592, -82.366592)
	manilla := routes.NewPoint("Manilla", 14.599512, 120.984222)
	// hongkong := routes.NewPoint("Hongkong", 22.302711, 114.177216)
	chicago := routes.NewPoint("Chicago", 41.881832, -87.623177)

	start := time.Now()
	arr := []routes.Point{amsterdam, kiruna, johannesburg, havana, chicago, goteborg, berlin, newyork, manilla}
	optimal := routes.OptimalPath(arr)
	fmt.Printf("\n%v cities", len(arr))
	fmt.Printf("\n%v permutations", optimal.NoOfPermutations)
	fmt.Printf("\nlowest score %v\n", optimal.Score)
	fmt.Printf("the optimal route is:")
	prettyPrintRoute(optimal.Points)
	elapsed := time.Since(start)
	fmt.Printf("calculated in %v milliseconds\n", elapsed.Milliseconds())
}
