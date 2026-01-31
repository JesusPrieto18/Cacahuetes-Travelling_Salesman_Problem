package tsp

import "math"

// NearestNeighbor constructs a tour using the nearest neighbor heuristic
// starting from the given city. Returns the tour and its total length.
func NearestNeighbor(inst *Instance, startCity int) ([]int, float64) {
	n := inst.NumCities
	visited := make([]bool, n)
	tour := make([]int, n)

	current := startCity
	tour[0] = current
	visited[current] = true

	for i := 1; i < n; i++ {
		nearest := -1
		nearestDist := math.MaxFloat64

		for j := 0; j < n; j++ {
			if !visited[j] && inst.Distance[current][j] < nearestDist {
				nearest = j
				nearestDist = inst.Distance[current][j]
			}
		}

		tour[i] = nearest
		visited[nearest] = true
		current = nearest
	}

	return tour, inst.TourLength(tour)
}

// BestNearestNeighbor tries nearest neighbor from all cities and returns the best tour
func BestNearestNeighbor(inst *Instance) ([]int, float64) {
	bestTour := make([]int, inst.NumCities)
	bestLength := math.MaxFloat64

	for start := 0; start < inst.NumCities; start++ {
		tour, length := NearestNeighbor(inst, start)
		if length < bestLength {
			bestLength = length
			copy(bestTour, tour)
		}
	}

	return bestTour, bestLength
}
