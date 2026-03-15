package geneticalgorithm

import "math/rand"

// TournamentSelection selects an individual using tournament selection of size k.
// Picks k random individuals and returns the one with lowest cost (best fitness).
func TournamentSelection(population []Individual, tournSize int) Individual {
	best := population[rand.Intn(len(population))]
	for i := 1; i < tournSize; i++ {
		candidate := population[rand.Intn(len(population))]
		if candidate.Cost < best.Cost {
			best = candidate
		}
	}
	return best
}
