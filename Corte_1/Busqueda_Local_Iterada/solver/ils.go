package solver

import (
	"math/rand"
	"tsp-ils/localsearch"
	"tsp-ils/models"
	"tsp-ils/perturbation"
	"tsp-ils/utils"
)

// Funcion busqueda local iterada
func ILS(ciudades []models.City, maxIteraciones int) ([]models.City, float64) {
	
	// Solución Inicial
	tourActual := utils.CopiarTour(ciudades)
	rand.Shuffle(len(tourActual), func(i, j int) {
		tourActual[i], tourActual[j] = tourActual[j], tourActual[i]
	})

	// Búsqueda Local Inicial
	tourActual, costoActual := localsearch.TwoOpt(tourActual)
	//fmt.Printf("   >> Costo Inicial (2-Opt puro): %.4f\n", costoActual)

	tourBest := utils.CopiarTour(tourActual)
	costoBest := costoActual

	// Bucle Principal
	for iter := 1; iter <= maxIteraciones; iter++ {
	
		// Perturbación
		tourCandidato := perturbation.DoubleBridge(tourActual)

		// Búsqueda Local
		tourCandidato, costoCandidato := localsearch.TwoOpt(tourCandidato)

		// Criterio de Aceptación
		if costoCandidato < costoActual {
			tourActual = tourCandidato
			costoActual = costoCandidato

			// Actualizar mejor global
			if costoActual < costoBest {
				tourBest = utils.CopiarTour(tourActual)
				costoBest = costoActual
				//fmt.Printf("   [Iter %d] ¡Nueva Mejor Solución! Costo: %.4f\n", iter, costoBest)
			}
		}
	}

	return tourBest, costoBest
}
