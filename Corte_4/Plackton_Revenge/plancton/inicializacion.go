package plancton

import (
	"math"
	"math/rand"
	"tsp/utils"
)

// FarthestInsertion construye un tour inicial de alta calidad basado en distancias máximas.
func FarthestInsertion(oceano utils.Oceano) []int {
	n := len(oceano)
	if n < 3 {
		perm := make([]int, n)
		for i := range perm {
			perm[i] = i
		}
		return perm
	}

	// 1. Iniciar con las dos ciudades más lejanas entre sí
	maxDist := -1.0
	var city1, city2 int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := utils.DistanciaEuclidiana(oceano[i], oceano[j])
			if d > maxDist {
				maxDist = d
				city1, city2 = i, j
			}
		}
	}

	// 2. Encontrar la tercera ciudad más lejana a estas dos
	maxMinDist := -1.0
	city3 := -1
	for i := 0; i < n; i++ {
		if i == city1 || i == city2 {
			continue
		}
		minD := math.Min(utils.DistanciaEuclidiana(oceano[city1], oceano[i]), 
						 utils.DistanciaEuclidiana(oceano[city2], oceano[i]))
		if minD > maxMinDist {
			maxMinDist = minD
			city3 = i
		}
	}

	tour := []int{city1, city2, city3}
	enTour := make([]bool, n)
	enTour[city1], enTour[city2], enTour[city3] = true, true, true

	// 3. Insertar el resto de las ciudades una por una
	for len(tour) < n {
		farthestCity := -1
		farthestDist := -1.0

		// Buscar la ciudad más lejana al tour actual
		for c := 0; c < n; c++ {
			if enTour[c] {
				continue
			}
			minDistAlTour := math.MaxFloat64
			for _, t := range tour {
				d := utils.DistanciaEuclidiana(oceano[c], oceano[t])
				if d < minDistAlTour {
					minDistAlTour = d
				}
			}
			if minDistAlTour > farthestDist {
				farthestDist = minDistAlTour
				farthestCity = c
			}
		}

		// Encontrar la mejor posición para insertar (mínimo incremento de costo)
		bestPos := -1
		minIncremento := math.MaxFloat64

		for pos := 0; pos < len(tour); pos++ {
			c1 := oceano[tour[pos]]
			c2 := oceano[tour[(pos+1)%len(tour)]]
			nueva := oceano[farthestCity]
			
			incremento := utils.DistanciaEuclidiana(c1, nueva) + 
						  utils.DistanciaEuclidiana(nueva, c2) - 
						  utils.DistanciaEuclidiana(c1, c2)
			
			if incremento < minIncremento {
				minIncremento = incremento
				bestPos = pos
			}
		}

		// Insertar la ciudad
		nuevoTour := make([]int, 0, len(tour)+1)
		nuevoTour = append(nuevoTour, tour[:bestPos+1]...)
		nuevoTour = append(nuevoTour, farthestCity)
		nuevoTour = append(nuevoTour, tour[bestPos+1:]...)
		tour = nuevoTour
		enTour[farthestCity] = true
	}

	return tour
}

// InicializarPoblacion crea la población inicial de planctones.
// Incluye una semilla de Farthest Insertion y el resto aleatorios para mantener diversidad.
func InicializarPoblacion(oceano utils.Oceano, nPop int) []Plancton {
	poblacion := make([]Plancton, 0, nPop)

	// 1. Crear el plancton "Alfa" con Farthest Insertion
	tourFI := FarthestInsertion(oceano)
	poblacion = append(poblacion, Plancton{
		Tour: tourFI,
		Cost: utils.CalcularCostoPermutacion(tourFI, oceano),
	})

	// 2. Llenar el resto de la población con permutaciones aleatorias
	nCities := len(oceano)
	for len(poblacion) < nPop {
		tourAleatorio := rand.Perm(nCities)
		poblacion = append(poblacion, Plancton{
			Tour: tourAleatorio,
			Cost: utils.CalcularCostoPermutacion(tourAleatorio, oceano),
		})
	}

	return poblacion
}