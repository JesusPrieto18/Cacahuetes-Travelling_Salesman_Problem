package geneticalgorithm

import (
	"math/rand"
	"sort"
)
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

// DoubleBridgeMutation realiza un 4-opt kick (saltos no secuenciales) para escapar de mínimos locales.
func DoubleBridgeMutation(tour []int) []int {
	n := len(tour)
	if n < 8 { return tour } // Requiere al menos 8 ciudades para funcionar bien

	// Elegir 4 puntos de corte aleatorios distintos
	cuts := []int{rand.Intn(n), rand.Intn(n), rand.Intn(n), rand.Intn(n)}
	sort.Ints(cuts)
	// Asegurar que sean únicos (simplificado para el ejemplo)
	for cuts[0] == cuts[1] || cuts[1] == cuts[2] || cuts[2] == cuts[3] {
		cuts = []int{rand.Intn(n), rand.Intn(n), rand.Intn(n), rand.Intn(n)}
		sort.Ints(cuts)
	}

	a, b, c, d := cuts[0], cuts[1], cuts[2], cuts[3]

	// Reensamblar en el orden: A-B, D-end, C-D, B-C (El cruce de puentes)
	newTour := make([]int, 0, n)
	newTour = append(newTour, tour[0:a]...)
	newTour = append(newTour, tour[c:d]...)
	newTour = append(newTour, tour[b:c]...)
	newTour = append(newTour, tour[a:b]...)
	newTour = append(newTour, tour[d:n]...)

	// Copiar de vuelta al tour original
	copy(tour, newTour)
	return tour
}