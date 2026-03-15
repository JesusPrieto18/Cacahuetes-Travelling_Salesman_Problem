package localsearch

import (
	"tsp-meme/models"
	"tsp-meme/utils"
)

// Funcion 2 opt para busqueda local (Adaptada para usar genotipos []int)
func TwoOpt(tour []int, cities []models.City) ([]int, float64) {
	// Copiar el tour para no modificar el original directamente
	mejorTour := make([]int, len(tour))
	copy(mejorTour, tour)
	
	// Calcular costo inicial (asumiendo que tienes EvaluateCost o iteramos)
	mejorCosto := calcularCosto(mejorTour, cities) 
	mejorado := true
	n := len(tour)

	for mejorado {
		mejorado = false
		for i := 1; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				// Buscar las distancias usando los índices del tour sobre el arreglo de ciudades
				d1 := utils.DistanciaEuclidiana(cities[mejorTour[i-1]], cities[mejorTour[i]])
				d2 := utils.DistanciaEuclidiana(cities[mejorTour[j]], cities[mejorTour[(j+1)%n]])
				costoActual := d1 + d2

				d3 := utils.DistanciaEuclidiana(cities[mejorTour[i-1]], cities[mejorTour[j]])
				d4 := utils.DistanciaEuclidiana(cities[mejorTour[i]], cities[mejorTour[(j+1)%n]])
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

// Funcion para invertir un segmento del tour basado en índices enteros
func invertirSegmento(tour []int, i, j int) {
	for i < j {
		tour[i], tour[j] = tour[j], tour[i]
		i++
		j--
	}
}

// Función auxiliar rápida para calcular el costo base en la búsqueda local
func calcularCosto(tour []int, cities []models.City) float64 {
	total := 0.0
	n := len(tour)
	for i := 0; i < n-1; i++ {
		total += utils.DistanciaEuclidiana(cities[tour[i]], cities[tour[i+1]])
	}
	total += utils.DistanciaEuclidiana(cities[tour[n-1]], cities[tour[0]])
	return total
}