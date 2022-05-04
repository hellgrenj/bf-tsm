package routes_test

import (
	"log"
	"os"
	"testing"

	"github.com/hellgrenj/bf-tsm/pkg/routes"
)

func quiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}
func TestOptimalPath(t *testing.T) {
	defer quiet()()
	amsterdam := routes.NewPoint("Amsterdam", 52.377956, 4.897070)
	berlin := routes.NewPoint("Berlin", 52.520008, 13.404954)
	kiruna := routes.NewPoint("Kiruna", 67.85000, 20.23000)
	goteborg := routes.NewPoint("Göteborg", 57.69000, 11.89000)
	johannesburg := routes.NewPoint("Johannes Burg", -26.195246, 28.034088)
	newyork := routes.NewPoint("New york", 40.730610, -73.935242)
	havana := routes.NewPoint("Havana", 23.113592, -82.366592)

	arr := []routes.Point{amsterdam, kiruna, johannesburg, havana, goteborg, berlin, newyork}
	optimal := routes.OptimalPath(arr)

	if optimal.Points[0].Label != "Amsterdam" {
		t.Errorf("Expected the first Point to be Amsterdam but it was %v", optimal.Points[0].Label)
	}
	if optimal.Points[1].Label != "Göteborg" {
		t.Errorf("Expected the second Point to be Göteborg but it was %v", optimal.Points[1].Label)
	}
	if optimal.Points[2].Label != "Kiruna" {
		t.Errorf("Expected the third Point to be Kiruna but it was %v", optimal.Points[2].Label)
	}
	if optimal.Points[3].Label != "Berlin" {
		t.Errorf("Expected the fourth Point to be Berlin but it was %v", optimal.Points[3].Label)
	}
	if optimal.Points[4].Label != "Johannes Burg" {
		t.Errorf("Expected the fifth Point to be Johannes Burg but it was %v", optimal.Points[4].Label)
	}
	if optimal.Points[5].Label != "Havana" {
		t.Errorf("Expected the sixth Point to be Havana but it was %v", optimal.Points[5].Label)
	}
	if optimal.Points[6].Label != "New york" {
		t.Errorf("Expected the last Point (before going back home) to be Havana but it was %v", optimal.Points[6].Label)
	}
	if optimal.NumberOfPermutations != 5040 { // 7 points or cities
		t.Errorf("Expected number of permuations to be 5040 but it was %v", optimal.NumberOfPermutations)
	}
}

func BenchmarkOptimalPath(b *testing.B) {
	defer quiet()()
	amsterdam := routes.NewPoint("Amsterdam", 52.377956, 4.897070)
	berlin := routes.NewPoint("Berlin", 52.520008, 13.404954)
	kiruna := routes.NewPoint("Kiruna", 67.85000, 20.23000)
	goteborg := routes.NewPoint("Göteborg", 57.69000, 11.89000)
	johannesburg := routes.NewPoint("Johannes Burg", -26.195246, 28.034088)
	newyork := routes.NewPoint("New york", 40.730610, -73.935242)

	for i := 0; i < b.N; i++ {
		arr := []routes.Point{amsterdam, kiruna, berlin, goteborg, johannesburg, newyork}
		routes.OptimalPath(arr)
	}
}
