package utils

import (
	"math"
	"tsp/models"
)

// Oceano es un alias para el mapa del problema. Representa el entorno 
// físico donde el plancton flota y busca nutrientes.
type Oceano = []models.City

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

// CalcularDistanciaAristas cuenta el número de aristas (conexiones) que difieren entre dos tours.
// Devuelve un valor entre 0 (idénticos) y N (totalmente diferentes).
func CalcularDistanciaAristas(tourA, tourB []int) int {
	if len(tourA) != len(tourB) || len(tourA) < 2 {
		return 0
	}
	n := len(tourA)

	// 1. Guardar todas las aristas del tourA en un mapa para búsqueda ultrarrápida (O(1))
	edgesA := make(map[[2]int]bool)
	for i := 0; i < n; i++ {
		u := tourA[i]
		v := tourA[(i+1)%n]
		// Ordenar siempre de menor a mayor para que la arista (5,8) sea igual a la (8,5)
		if u > v {
			u, v = v, u
		}
		edgesA[[2]int{u, v}] = true
	}

	distancia := 0

	// 2. Recorrer el tourB y ver cuáles de sus aristas NO existen en el tourA
	for i := 0; i < n; i++ {
		u := tourB[i]
		v := tourB[(i+1)%n]
		if u > v {
			u, v = v, u
		}
		
		// Si la arista del tourB no está en el mapa de tourA, sumamos a la distancia
		if !edgesA[[2]int{u, v}] {
			distancia++
		}
	}

	return distancia
}

// CalcularCostoPermutacion evalúa el costo total de una ruta representada por índices.
// (Equivalente a EvaluateCost pero independiente del paquete geneticalgorithm).
func CalcularCostoPermutacion(tour []int, cities []models.City) float64 {
	total := 0.0
	n := len(tour)
	for i := 0; i < n-1; i++ {
		total += DistanciaEuclidiana(cities[tour[i]], cities[tour[i+1]])
	}
	total += DistanciaEuclidiana(cities[tour[n-1]], cities[tour[0]])
	return total
}

// CopiarPermutacion crea una copia independiente de un slice de enteros.
func CopiarPermutacion(tour []int) []int {
	c := make([]int, len(tour))
	copy(c, tour)
	return c
}