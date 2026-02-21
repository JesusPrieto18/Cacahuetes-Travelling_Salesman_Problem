package solver

import (
	"tsp-common/models"
	"tsp-sa/simulatedannealing"
)

// SimulatedAnnealingSolver recibe un tour inicial (que puede venir de Local Search)
// y lo mejora usando Recocido Simulado.
func SimulatedAnnealingSolver(tourInicial []models.City, costoInicial float64, config simulatedannealing.SAConfig) ([]models.City, float64) {
	
	// Ejecutar SA
	mejorTour, mejorCosto := simulatedannealing.EjecutarSA(tourInicial, config)

	return mejorTour, mejorCosto
}