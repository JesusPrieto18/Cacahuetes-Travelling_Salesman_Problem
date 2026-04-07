package plancton

import (
	"math/rand"
	"tsp/models"
	"tsp/utils"
)

// AplicarTurbulencia dispersa aleatoriamente una fracción (mu) de la población.
// Modifica la población en su lugar y protege al individuo en el índice 0.
func AplicarTurbulencia(poblacion []Plancton, mu float64, cities []models.City) {
	nPop := len(poblacion)
	turbCount := int(mu * float64(nPop))

	// Si la población es muy pequeña o la intensidad es 0, omitimos la turbulencia
	if turbCount == 0 || nPop < 2 {
		return
	}

	nCities := len(cities)

	for i := 0; i < turbCount; i++ {
		// Elegir un índice aleatorio evitando el 0.
		// rand.Intn(nPop-1) genera [0, nPop-2], al sumar 1 obtenemos [1, nPop-1].
		// Esto protege al campeón si asumimos que la población fue ordenada previamente.
		idx := rand.Intn(nPop-1) + 1

		// Reinicialización total usando rand.Perm nativo de Go
		nuevoTour := rand.Perm(nCities)

		// Actualizar el plancton arrastrado por la turbulencia
		poblacion[idx].Tour = nuevoTour
		poblacion[idx].Cost = utils.CalcularCostoPermutacion(nuevoTour, cities)
	}
}