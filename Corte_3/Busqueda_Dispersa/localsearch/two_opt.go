package localsearch

import (
	"tsp-ds/models"
	"tsp-ds/utils"
)

// Funcion 2 opt para busqueda local (Adaptada y protegida contra bucles)
func TwoOpt(tour []int, cities []models.City) ([]int, float64) {
	mejorTour := make([]int, len(tour))
	copy(mejorTour, tour)

	mejorado := true
	n := len(tour)
	for mejorado {
		mejorado = false
		for i := 1; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				d1 := utils.DistanciaEuclidiana(cities[mejorTour[i-1]], cities[mejorTour[i]])
				d2 := utils.DistanciaEuclidiana(cities[mejorTour[j]], cities[mejorTour[(j+1)%n]])
				costoActual := d1 + d2

				d3 := utils.DistanciaEuclidiana(cities[mejorTour[i-1]], cities[mejorTour[j]])
				d4 := utils.DistanciaEuclidiana(cities[mejorTour[i]], cities[mejorTour[(j+1)%n]])
				costoNuevo := d3 + d4

				// EL ARREGLO ESTÁ AQUÍ: Añadimos un margen de tolerancia (0.0001)
				if (costoActual - costoNuevo) > 0.0001 {
					invertirSegmento(mejorTour, i, j)
					mejorado = true
				}
			}
		}
	}
	
	// Calculamos el costo final directamente del tour terminado para evitar
	// arrastrar errores de precisión acumulados en restas anteriores.
	mejorCostoFinal := calcularCosto(mejorTour, cities)
	
	return mejorTour, mejorCostoFinal
}

func invertirSegmento(tour []int, i, j int) {
	for i < j {
		tour[i], tour[j] = tour[j], tour[i]
		i++
		j--
	}
}

func calcularCosto(tour []int, cities []models.City) float64 {
	total := 0.0
	n := len(tour)
	for i := 0; i < n-1; i++ {
		total += utils.DistanciaEuclidiana(cities[tour[i]], cities[tour[i+1]])
	}
	total += utils.DistanciaEuclidiana(cities[tour[n-1]], cities[tour[0]])
	return total
}