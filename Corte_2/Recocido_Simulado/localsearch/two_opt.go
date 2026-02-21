package localsearch

import (
	"tsp-common/models"
	"tsp-common/utils"
)

// Funcion 2 opt para busqueda local
func TwoOpt(tour []models.City) ([]models.City, float64) {
	mejorTour := utils.CopiarTour(tour)
	mejorCosto := utils.CalcularCostoTotal(mejorTour)
	mejorado := true
	n := len(tour)

	for mejorado {
		mejorado = false
		for i := 1; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				d1 := utils.DistanciaEuclidiana(mejorTour[i-1], mejorTour[i])
				d2 := utils.DistanciaEuclidiana(mejorTour[j], mejorTour[(j+1)%n])
				costoActual := d1 + d2

				d3 := utils.DistanciaEuclidiana(mejorTour[i-1], mejorTour[j])
				d4 := utils.DistanciaEuclidiana(mejorTour[i], mejorTour[(j+1)%n])
				costoNuevo := d3 + d4

				if costoNuevo < costoActual {
					invertirSegmento(mejorTour, i, j)
					mejorCosto -= (costoActual - costoNuevo)
					mejorado = true
				}
			}
		}
	}
	return mejorTour, mejorCosto
}

// Funcion para invertir un segmento del tour
func invertirSegmento(tour []models.City, i, j int) {
	for i < j {
		tour[i], tour[j] = tour[j], tour[i]
		i++
		j--
	}
}
