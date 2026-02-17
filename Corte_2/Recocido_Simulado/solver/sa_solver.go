package solver

import (
	"tsp-ils/models"
	"tsp-ils/simulatedannealing"
)

// SimulatedAnnealingSolver recibe un tour inicial (que puede venir de Local Search)
// y lo mejora usando Recocido Simulado.
func SimulatedAnnealingSolver(tourInicial []models.City, costoInicial float64) ([]models.City, float64) {
	
	// Configuración de parámetros
	// NOTA: T0 debe ser suficientemente alta para permitir salir del óptimo local 
	// que nos entregó la búsqueda local.
	config := simulatedannealing.SAConfig{
		InitialTemp: 1000.0,  // Temperatura inicial
		Alpha:       0.995,   // Enfriamiento lento para mayor calidad
		MinTemp:     0.001,   // Criterio de parada
		IterPerTemp: 1000,    // Intensificación por nivel
	}

	// Ejecutar SA
	mejorTour, mejorCosto := simulatedannealing.EjecutarSA(tourInicial, config)

	return mejorTour, mejorCosto
}