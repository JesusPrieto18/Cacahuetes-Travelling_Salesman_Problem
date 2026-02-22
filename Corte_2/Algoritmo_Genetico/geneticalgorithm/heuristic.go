package geneticalgorithm

import (
	"math"
	"tsp-ga/models"
	"tsp-ga/utils"
)

// FarthestInsertion builds a tour using the farthest insertion heuristic.
// Adapted from Corte_1/Heuristica/tsp/insertion.go to work with []models.City.
// Returns a permutation of indices [0..n-1].
func FarthestInsertion(cities []models.City) []int {
	n := len(cities)
	if n < 3 {
		perm := make([]int, n)
		for i := range perm {
			perm[i] = i
		}
		return perm
	}

	// Precompute distance matrix for efficiency (O(n^2) memory, avoids repeated sqrt)
	dist := make([][]float64, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := utils.DistanciaEuclidiana(cities[i], cities[j])
			dist[i][j] = d
			dist[j][i] = d
		}
	}

	// Start with the two farthest cities
	maxDist := -1.0
	var city1, city2 int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if dist[i][j] > maxDist {
				maxDist = dist[i][j]
				city1, city2 = i, j
			}
		}
	}

	// Find the city farthest from these two (max of min distances)
	maxMinDist := -1.0
	city3 := -1
	for i := 0; i < n; i++ {
		if i == city1 || i == city2 {
			continue
		}
		minDist := math.Min(dist[city1][i], dist[city2][i])
		if minDist > maxMinDist {
			maxMinDist = minDist
			city3 = i
		}
	}

	// Initialize tour with these 3 cities
	tour := []int{city1, city2, city3}
	inTour := make([]bool, n)
	inTour[city1] = true
	inTour[city2] = true
	inTour[city3] = true

	// Insert remaining cities one by one
	for len(tour) < n {
		// Find the city farthest from the current tour
		farthestCity := -1
		farthestDist := -1.0

		for c := 0; c < n; c++ {
			if inTour[c] {
				continue
			}
			// Minimum distance from c to any city in the tour
			minDist := math.MaxFloat64
			for _, t := range tour {
				if dist[c][t] < minDist {
					minDist = dist[c][t]
				}
			}
			if minDist > farthestDist {
				farthestDist = minDist
				farthestCity = c
			}
		}

		// Find the best position to insert (minimizes cost increase)
		bestPos := -1
		bestCost := math.MaxFloat64

		for pos := 0; pos < len(tour); pos++ {
			i := tour[pos]
			j := tour[(pos+1)%len(tour)]
			costIncrease := dist[i][farthestCity] + dist[farthestCity][j] - dist[i][j]
			if costIncrease < bestCost {
				bestCost = costIncrease
				bestPos = pos
			}
		}

		// Insert the city after bestPos
		newTour := make([]int, len(tour)+1)
		copy(newTour, tour[:bestPos+1])
		newTour[bestPos+1] = farthestCity
		copy(newTour[bestPos+2:], tour[bestPos+1:])
		tour = newTour
		inTour[farthestCity] = true
	}

	return tour
}
