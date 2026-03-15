package geneticalgorithm

import (
	"math/rand"
	"sort"
	"tsp-ds/models"
	"tsp-ds/utils"
	"tsp-ds/localsearch"
)

// GAConfig holds the genetic algorithm parameters.
type GAConfig struct {
	PopSize        int     // Population size
	Generations    int     // Maximum number of generations
	MutationRate   float64 // Mutation probability
	TournamentSize int     // Tournament size for selection
	StagnationLimit int    // Stop after this many generations without improvement (0 = disabled)
	RelinkPct       float64 // NUEVO: % de pares a reenlazar (ej. 0.5 para 50%)
	DivThreshold    int     // NUEVO: Distancia mínima (aristas) para aceptar un individuo (ej. 5)
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

	// 2. Bucle Generacional (Scatter Search)
	for gen := 0; gen < config.Generations; gen++ {
		totalGens = gen + 1
		offspring := make([]Individual, 0)

		// A. COMBINACIÓN SISTEMÁTICA (Path Relinking para un % de los pares - Inciso b)
		nPop := len(population)
		for i := 0; i < nPop-1; i++ {
			for j := i + 1; j < nPop; j++ {
				// Solo procesamos un porcentaje dado de todos los pares posibles
				if rand.Float64() < config.RelinkPct {
					
					// Reenlazado de Caminos (Path Relinking) entre el individuo i y el individuo j
					childTour := PathRelinking(population[i].Tour, population[j].Tour, cities)

					// Mutación (Opcional, para evitar estancamiento total)
					if rand.Float64() < config.MutationRate {
						childTour = DoubleBridgeMutation(childTour)
					}

					// Búsqueda Local
					childTourOpt, childCost := localsearch.TwoOpt(childTour, cities)

					offspring = append(offspring, Individual{Tour: childTourOpt, Cost: childCost})
				}
			}
		}

		// Combinar la población actual con las nuevas soluciones generadas
		combined := make([]Individual, 0, len(population)+len(offspring))
		combined = append(combined, population...)
		combined = append(combined, offspring...)

		// Ordenar todo por costo (los mejores primero)
		sort.Slice(combined, func(i, j int) bool {
			return combined[i].Cost < combined[j].Cost
		})

		// B. ACTUALIZACIÓN DEL CONJUNTO DE REFERENCIA CON DIVERSIDAD (Inciso a)
		newPop := make([]Individual, 0, config.PopSize)
		rejected := make([]Individual, 0) // Guardamos los rechazados por si nos faltan individuos

		// El campeón siempre pasa a la siguiente generación
		newPop = append(newPop, combined[0])

		for i := 1; i < len(combined); i++ {
			candidate := combined[i]
			esDiverso := true

			// Medir la distancia contra los que YA entraron a la nueva población
			for _, selected := range newPop {
				distancia := utils.CalcularDistanciaAristas(candidate.Tour, selected.Tour)
				
				// Si comparte demasiadas aristas (la distancia es muy baja), lo rechazamos
				if distancia < config.DivThreshold {
					esDiverso = false
					break
				}
			}

			if esDiverso {
				newPop = append(newPop, candidate)
			} else {
				rejected = append(rejected, candidate)
			}

			// Si ya llenamos el Reference Set, dejamos de buscar
			if len(newPop) == config.PopSize {
				break
			}
		}

		// Si el filtro de diversidad fue TAN estricto que no logramos llenar la población (PopSize),
		// rellenamos los espacios vacíos con los mejores rechazados.
		for i := 0; len(newPop) < config.PopSize && i < len(rejected); i++ {
			newPop = append(newPop, rejected[i])
		}

		population = newPop

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
