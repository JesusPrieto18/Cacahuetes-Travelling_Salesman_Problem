package acs

import (
	"heuristica/tsp"
	"math"
	"math/rand"
)

// ACS implements the Ant Colony System algorithm
type ACS struct {
	instance   *tsp.Instance
	params     Params
	pheromone  *PheromoneMatrix
	ants       []*Ant
	eta        [][]float64 // heuristic information (1/distance)
	bestTour   []int
	bestLength float64
	rng        *rand.Rand
}

// New creates a new ACS solver for the given TSP instance
func New(instance *tsp.Instance, params Params, seed int64) *ACS {
	n := instance.NumCities

	// Calculate nearest neighbor tour length for tau0
	_, lnn := tsp.NearestNeighbor(instance, 0)

	// Initialize pheeromone matrix
	pheromone := NewPheromoneMatrix(n, lnn)

	// Create ants
	ants := make([]*Ant, params.NumAnts)
	for i := range ants {
		ants[i] = NewAnt(n)
	}

	// Calculate heuristic informatcion
	eta := make([][]float64, n)
	for i := range eta {
		eta[i] = make([]float64, n)
		for j := range eta[i] {
			if i != j && instance.Distance[i][j] > 0 {
				eta[i][j] = 1.0 / instance.Distance[i][j]
			}
		}
	}

	return &ACS{
		instance:   instance,
		params:     params,
		pheromone:  pheromone,
		ants:       ants,
		eta:        eta,
		bestTour:   nil,
		bestLength: math.MaxFloat64,
		rng:        rand.New(rand.NewSource(seed)),
	}
}

// Run executes the ACS algorithm and returns the best tour found
func (acs *ACS) Run() ([]int, float64) {
	for iter := 0; iter < acs.params.Iterations; iter++ {
		for _, ant := range acs.ants {
			tour := acs.constructSolution(ant)
			length := acs.instance.TourLength(tour)

			if length < acs.bestLength {
				acs.bestLength = length
				acs.bestTour = make([]int, len(tour))
				copy(acs.bestTour, tour)
			}
		}

		acs.pheromone.GlobalUpdate(acs.bestTour, acs.bestLength, acs.params.Alpha)
	}

	return acs.bestTour, acs.bestLength
}

// constructSolution builds a complete tour for an ant
func (acs *ACS) constructSolution(ant *Ant) []int {
	n := acs.instance.NumCities

	// Start from a random city
	startCity := acs.rng.Intn(n)
	ant.Reset(startCity)

	// Build tour city by city
	for !ant.TourComplete() {
		current := ant.CurrentCity()
		next := acs.selectNextCity(ant, current)
		ant.Visit(next)

		// Local pheromone update
		acs.pheromone.LocalUpdate(current, next, acs.params.Rho)
	}

	return ant.Tour()
}

// selectNextCity implements the ACS state transition rule
func (acs *ACS) selectNextCity(ant *Ant, current int) int {
	q := acs.rng.Float64()

	if q <= acs.params.Q0 {
		// Exploitation: choose best city deterministically
		return acs.exploitationChoice(ant, current)
	}
	// Exploration: probabilistic selection
	return acs.explorationChoice(ant, current)
}

// exploitationChoice returns argmax
func (acs *ACS) exploitationChoice(ant *Ant, current int) int {
	bestCity := -1
	bestValue := -1.0

	for j := 0; j < acs.instance.NumCities; j++ {
		if ant.CanVisit(j) {
			tau := acs.pheromone.Get(current, j)
			eta := acs.eta[current][j]
			value := tau * math.Pow(eta, acs.params.Beta)

			if value > bestValue {
				bestValue = value
				bestCity = j
			}
		}
	}

	return bestCity
}

// explorationChoice uses roulette wheel selection
func (acs *ACS) explorationChoice(ant *Ant, current int) int {
	unvisited := ant.UnvisitedCities()
	if len(unvisited) == 1 {
		return unvisited[0]
	}

	// Calculate probabilities
	probs := make([]float64, len(unvisited))
	total := 0.0

	for i, j := range unvisited {
		tau := acs.pheromone.Get(current, j)
		eta := acs.eta[current][j]
		probs[i] = tau * math.Pow(eta, acs.params.Beta)
		total += probs[i]
	}

	// Normalize
	for i := range probs {
		probs[i] /= total
	}

	// Roulette wheel selection
	r := acs.rng.Float64()
	cumulative := 0.0

	for i, p := range probs {
		cumulative += p
		if r <= cumulative {
			return unvisited[i]
		}
	}

	// Fallback (shouldn't happen)
	return unvisited[len(unvisited)-1]
}

// BestTour returns the best tour found so far
func (acs *ACS) BestTour() []int {
	return acs.bestTour
}

// BestLength returns the length of the best tour found so far
func (acs *ACS) BestLength() float64 {
	return acs.bestLength
}
