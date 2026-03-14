package geneticalgorithm

import "math/rand"

// InversionMutation implements inversion mutation - Clase 7, slide 18.
// Picks 2 random positions and reverses the segment between them.
func InversionMutation(tour []int) {
	n := len(tour)
	i := rand.Intn(n)
	j := rand.Intn(n)
	if i > j {
		i, j = j, i
	}
	for i < j {
		tour[i], tour[j] = tour[j], tour[i]
		i++
		j--
	}
}
