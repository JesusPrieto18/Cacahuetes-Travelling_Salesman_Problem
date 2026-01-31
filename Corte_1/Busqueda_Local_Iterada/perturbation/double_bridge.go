package perturbation

import (
	"math/rand"
	"tsp-ils/models"
	"tsp-ils/utils"
)

// Funcion Double Bridge para la perturbacion
func DoubleBridge(tour []models.City) []models.City {
	n := len(tour)
	if n < 8 {
		return utils.CopiarTour(tour)
	}

	indices := make([]int, 3)
	indices[0] = rand.Intn(n/4) + 1
	indices[1] = indices[0] + rand.Intn(n/4) + 1
	indices[2] = indices[1] + rand.Intn(n/4) + 1
	if indices[2] >= n {
		indices[2] = n - 1
	}

	pos1, pos2, pos3 := indices[0], indices[1], indices[2]

	A := tour[:pos1]
	B := tour[pos1:pos2]
	C := tour[pos2:pos3]
	D := tour[pos3:]

	nuevoTour := make([]models.City, 0, n)
	nuevoTour = append(nuevoTour, A...)
	nuevoTour = append(nuevoTour, D...)
	nuevoTour = append(nuevoTour, C...)
	nuevoTour = append(nuevoTour, B...)

	return nuevoTour
}
