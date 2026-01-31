package tsp

// Instance represents a TSP problem instance
type Instance struct {
	NumCities   int
	Distance    [][]float64
	OptimalCost float64
}

// TourLength calculates the total length of a tour
func (inst *Instance) TourLength(tour []int) float64 {
	total := 0.0
	n := len(tour)
	for i := 0; i < n; i++ {
		from := tour[i]
		to := tour[(i+1)%n]
		total += inst.Distance[from][to]
	}
	return total
}
