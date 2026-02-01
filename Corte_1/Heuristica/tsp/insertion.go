package tsp

import "math"

// FarthestInsertion constructs a tour using the farthest insertion heuristic
func FarthestInsertion(inst *Instance) ([]int, float64) {
	n := inst.NumCities
	if n < 3 {
		return nil, 0
	}

	// Start with the two farthest cities
	maxDist := -1.0
	var city1, city2 int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if inst.Distance[i][j] > maxDist {
				maxDist = inst.Distance[i][j]
				city1, city2 = i, j
			}
		}
	}

	// Find the city farthest from these two
	maxMinDist := -1.0
	city3 := -1
	for i := 0; i < n; i++ {
		if i == city1 || i == city2 {
			continue
		}
		minDist := math.Min(inst.Distance[city1][i], inst.Distance[city2][i])
		if minDist > maxMinDist {
			maxMinDist = minDist
			city3 = i
		}
	}

	// Initialize tour
	tour := []int{city1, city2, city3}
	inTour := make([]bool, n)
	inTour[city1] = true
	inTour[city2] = true
	inTour[city3] = true

	// Insert remaining cities one by one
	for len(tour) < n {
		// Find the city farthest from the tour
		farthestCity := -1
		maxMinDist := -1.0

		for c := 0; c < n; c++ {
			if inTour[c] {
				continue
			}

			// Find minimum distance from c to any city in tour
			minDist := math.MaxFloat64
			for _, t := range tour {
				if inst.Distance[c][t] < minDist {
					minDist = inst.Distance[c][t]
				}
			}

			if minDist > maxMinDist {
				maxMinDist = minDist
				farthestCity = c
			}
		}

		// Find the best position to insert this city
		bestPos := -1
		bestCost := math.MaxFloat64

		for pos := 0; pos < len(tour); pos++ {
			i := tour[pos]
			j := tour[(pos+1)%len(tour)]
			costIncrease := inst.Distance[i][farthestCity] + inst.Distance[farthestCity][j] - inst.Distance[i][j]

			if costIncrease < bestCost {
				bestCost = costIncrease
				bestPos = pos
			}
		}

		// Insert the city
		newTour := make([]int, len(tour)+1)
		copy(newTour, tour[:bestPos+1])
		newTour[bestPos+1] = farthestCity
		copy(newTour[bestPos+2:], tour[bestPos+1:])
		tour = newTour
		inTour[farthestCity] = true
	}

	return tour, inst.TourLength(tour)
}
