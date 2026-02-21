package grasp

import (
	"math/rand"
	"tsp-sa/localsearch"
	"tsp-sa/models"
	"tsp-sa/utils"
)

func GraspReactivo(cities []models.City, maxIter int) ([]models.City, float64) {
	var bestTour []models.City
	bestCost := 1e18

	// Inicializar opciones de alpha
	alphas := []*AlphaOption{
		{value: 0.1}, {value: 0.2}, {value: 0.3}, {value: 0.4},
	}

	for i := 1; i <= maxIter; i++ {
		// Seleccionar un alpha al azar
		alphaIdx := rand.Intn(len(alphas))
		alphaOpt := alphas[alphaIdx]

		// Construccion con punto de inicio aleatorio
		initialSolution := buildSolution(cities, alphaOpt.value)

		// Busqueda local
		refinedTour, refinedCost := localsearch.TwoOpt(initialSolution)

		alphaOpt.costSum += refinedCost
		alphaOpt.uses++

		if refinedCost < bestCost {
			bestCost = refinedCost
			bestTour = utils.CopiarTour(refinedTour)
		}
	}

	return bestTour, bestCost
}
