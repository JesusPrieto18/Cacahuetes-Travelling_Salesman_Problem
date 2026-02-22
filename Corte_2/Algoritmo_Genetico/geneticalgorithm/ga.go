package geneticalgorithm

import (
	"math/rand"
	"sort"
	"tsp-ga/models"
	"tsp-ga/utils"
)

// GAConfig holds the genetic algorithm parameters.
type GAConfig struct {
	PopSize        int     // Population size
	Generations    int     // Maximum number of generations
	MutationRate   float64 // Mutation probability
	TournamentSize int     // Tournament size for selection
	StagnationLimit int    // Stop after this many generations without improvement (0 = disabled)
}

// Individual represents a candidate solution (genotype: index permutation).
type Individual struct {
	Tour []int   // Genotype: permutation of indices [0..N-1]
	Cost float64 // Fitness: total tour cost
}

// EvaluateCost computes the cost of a tour given as an index permutation.
func EvaluateCost(tour []int, cities []models.City) float64 {
	total := 0.0
	n := len(tour)
	for i := 0; i < n-1; i++ {
		total += utils.DistanciaEuclidiana(cities[tour[i]], cities[tour[i+1]])
	}
	total += utils.DistanciaEuclidiana(cities[tour[n-1]], cities[tour[0]])
	return total
}

// randomPermutation creates a random permutation of [0..n-1].
func randomPermutation(n int) []int {
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}
	rand.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	return perm
}

// copyTour copies an int slice.
func copyTour(tour []int) []int {
	c := make([]int, len(tour))
	copy(c, tour)
	return c
}

// perturbTour applies k random swaps to a copy of the tour.
func perturbTour(tour []int, k int) []int {
	p := copyTour(tour)
	n := len(p)
	for s := 0; s < k; s++ {
		i := rand.Intn(n)
		j := rand.Intn(n)
		p[i], p[j] = p[j], p[i]
	}
	return p
}

// isDuplicate checks if a cost already exists in the population (same cost = likely duplicate).
func isDuplicate(pop []Individual, cost float64) bool {
	for _, ind := range pop {
		if ind.Cost == cost {
			return true
		}
	}
	return false
}

// initPopulation creates the initial population with:
//   - 1 Farthest Insertion individual
//   - ~15% perturbed variants of the FI tour
//   - ~85% random permutations
//   - Duplicate costs are discarded and regenerated.
func initPopulation(cities []models.City, popSize int) []Individual {
	n := len(cities)
	pop := make([]Individual, 0, popSize)

	// 1. Farthest Insertion seed
	fiTour := FarthestInsertion(cities)
	fiCost := EvaluateCost(fiTour, cities)
	pop = append(pop, Individual{Tour: fiTour, Cost: fiCost})

	// 2. Perturbed variants of FI tour (~15% of population)
	numPerturbed := popSize * 15 / 100
	swaps := n / 5 // number of swaps per perturbation
	if swaps < 2 {
		swaps = 2
	}
	for i := 0; i < numPerturbed; i++ {
		pt := perturbTour(fiTour, swaps)
		cost := EvaluateCost(pt, cities)
		if !isDuplicate(pop, cost) {
			pop = append(pop, Individual{Tour: pt, Cost: cost})
		}
	}

	// 3. Random permutations for the rest, with diversity control
	maxAttempts := popSize * 3 // avoid infinite loop
	attempts := 0
	for len(pop) < popSize && attempts < maxAttempts {
		tour := randomPermutation(n)
		cost := EvaluateCost(tour, cities)
		if !isDuplicate(pop, cost) {
			pop = append(pop, Individual{Tour: tour, Cost: cost})
		}
		attempts++
	}

	// If we still need more (very unlikely), fill without diversity check
	for len(pop) < popSize {
		tour := randomPermutation(n)
		pop = append(pop, Individual{Tour: tour, Cost: EvaluateCost(tour, cities)})
	}

	return pop
}

// GAResult holds the output of a GA run including convergence info.
type GAResult struct {
	BestTour     []models.City
	BestCost     float64
	LastImproveGen int    // Generation where the last improvement occurred
	TotalGens    int      // Total generations executed
	StopReason   string   // "max_generaciones" or "estancamiento"
}

// RunGA executes the genetic algorithm and returns the result with convergence info.
func RunGA(cities []models.City, config GAConfig) GAResult {
	n := len(cities)

	// 1. Initialize diverse population
	population := initPopulation(cities, config.PopSize)

	// Find initial best
	best := population[0]
	for _, ind := range population[1:] {
		if ind.Cost < best.Cost {
			best = Individual{Tour: copyTour(ind.Tour), Cost: ind.Cost}
		}
	}

	// Convergence tracking
	stagnationCount := 0
	lastImproveGen := 0
	totalGens := 0
	stopReason := "max_generaciones"

	// 2. Generational loop
	for gen := 0; gen < config.Generations; gen++ {
		totalGens = gen + 1

		// Generate offspring (λ = PopSize)
		offspring := make([]Individual, 0, config.PopSize)

		for len(offspring) < config.PopSize {
			// Select parents by tournament
			parent1 := TournamentSelection(population, config.TournamentSize)
			parent2 := TournamentSelection(population, config.TournamentSize)

			// Cut and Fill crossover
			child1Tour, child2Tour := CutAndFillCrossover(parent1.Tour, parent2.Tour)

			// Inversion mutation with probability MutationRate
			if rand.Float64() < config.MutationRate {
				InversionMutation(child1Tour)
			}
			if rand.Float64() < config.MutationRate {
				InversionMutation(child2Tour)
			}

			// Evaluate offspring
			child1 := Individual{Tour: child1Tour, Cost: EvaluateCost(child1Tour, cities)}
			child2 := Individual{Tour: child2Tour, Cost: EvaluateCost(child2Tour, cities)}

			offspring = append(offspring, child1)
			if len(offspring) < config.PopSize {
				offspring = append(offspring, child2)
			}
		}

		// (μ+λ) survivor selection: merge population + offspring, keep best PopSize
		combined := make([]Individual, 0, len(population)+len(offspring))
		combined = append(combined, population...)
		combined = append(combined, offspring...)

		sort.Slice(combined, func(i, j int) bool {
			return combined[i].Cost < combined[j].Cost
		})

		population = combined[:config.PopSize]

		// Update global best
		if population[0].Cost < best.Cost {
			best = Individual{Tour: copyTour(population[0].Tour), Cost: population[0].Cost}
			stagnationCount = 0
			lastImproveGen = gen + 1
		} else {
			stagnationCount++
		}

		// Stagnation termination
		if config.StagnationLimit > 0 && stagnationCount >= config.StagnationLimit {
			stopReason = "estancamiento"
			break
		}
	}

	// 3. Convert best tour (indices) to []City
	bestTourCities := make([]models.City, n)
	for i, idx := range best.Tour {
		bestTourCities[i] = cities[idx]
	}

	return GAResult{
		BestTour:       bestTourCities,
		BestCost:       best.Cost,
		LastImproveGen: lastImproveGen,
		TotalGens:      totalGens,
		StopReason:     stopReason,
	}
}
