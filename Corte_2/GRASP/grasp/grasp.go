package grasp

import (
	"tsp-sa/localsearch"
	"tsp-sa/models"
	"tsp-sa/utils"
)

func GraspReactivo(cities []models.City, maxIter int, alpha *Alpha) ([]models.City, float64) {
	var bestTour []models.City
	bestCost := 1e18


	for i := 1; i <= maxIter; i++ {
		// Construccion con punto de inicio aleatorio
		initialSolution := buildSolution(cities, alpha.Value)

		// Busqueda local
		refinedTour, refinedCost := localsearch.TwoOpt(initialSolution)

		alpha.costSum += refinedCost
		alpha.uses++

		if refinedCost < bestCost {
			bestCost = refinedCost
			bestTour = utils.CopiarTour(refinedTour)
		}
	}

	return bestTour, bestCost
}
