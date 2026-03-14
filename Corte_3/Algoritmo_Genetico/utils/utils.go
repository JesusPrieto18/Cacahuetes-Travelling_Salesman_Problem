package utils

import (
	"math"
	"tsp-ga/models"
)

// Funcion para calcular la distancia euclidiana entre dos ciudades
func DistanciaEuclidiana(c1, c2 models.City) float64 {
	return math.Sqrt(math.Pow(c1.X-c2.X, 2) + math.Pow(c1.Y-c2.Y, 2))
}

// Funcion para calcular el costo total de un tour
func CalcularCostoTotal(tour []models.City) float64 {
	total := 0.0
	for i := 0; i < len(tour)-1; i++ {
		total += DistanciaEuclidiana(tour[i], tour[i+1])
	}
	total += DistanciaEuclidiana(tour[len(tour)-1], tour[0])
	return total
}

// Funcion para copiar un tour
func CopiarTour(tour []models.City) []models.City {
	nueva := make([]models.City, len(tour))
	copy(nueva, tour)
	return nueva
}
