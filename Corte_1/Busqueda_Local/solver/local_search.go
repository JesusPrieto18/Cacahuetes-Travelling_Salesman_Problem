package solver

import (
	"fmt"
	"math/rand"
	"tsp-ils/localsearch"
	"tsp-ils/models"
	"tsp-ils/utils"
)

// LocalSearch ejecuta el algoritmo de Búsqueda
// Genera un inicio aleatorio y aplica 2-opt hasta llegar a un óptimo local.
func LocalSearch(ciudades []models.City) ([]models.City, float64) {
	
	// Solución Inicial Aleatoria
	tourActual := utils.CopiarTour(ciudades)
	
	// Aleatorizamos el orden (Random Start)
	rand.Shuffle(len(tourActual), func(i, j int) {
		tourActual[i], tourActual[j] = tourActual[j], tourActual[i]
	})

	costoInicial := utils.CalcularCostoTotal(tourActual)
	fmt.Printf("   >> Costo Inicial (Aleatorio): %.4f\n", costoInicial)

	// Aplicar 2-Opt
	mejorTour, mejorCosto := localsearch.TwoOpt(tourActual)

	return mejorTour, mejorCosto
}
