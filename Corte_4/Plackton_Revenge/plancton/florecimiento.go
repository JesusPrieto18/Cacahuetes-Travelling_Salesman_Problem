package plancton

import (
	"math/rand"
	"sort"
	"tsp/models"
	"tsp/utils"
)

// AplicarFlorecimiento toma la fracción de élite (B) y genera descendientes
// con una perturbación mínima. Retorna la población incrementada.
func AplicarFlorecimiento(poblacion []Plancton, B float64, cities []models.City) []Plancton {
	nPop := len(poblacion)
	if nPop == 0 {
		return poblacion
	}

	// 1. Ordenar la población para asegurar que los primeros elementos son la élite
	sort.Slice(poblacion, func(i, j int) bool {
		return poblacion[i].Cost < poblacion[j].Cost
	})

	// Calcular la cantidad de élites a reproducir
	eliteCount := int(B * float64(nPop))
	if eliteCount < 1 {
		eliteCount = 1 // Garantizar al menos un florecimiento si B > 0
	}

	nCities := len(cities)
	hijos := make([]Plancton, 0, eliteCount)

	// 2. Reproducción de la élite
	for i := 0; i < eliteCount; i++ {
		padre := poblacion[i]
		tourHijo := utils.CopiarPermutacion(padre.Tour)

		// Perturbación topológica (Inversión de sub-ruta en lugar de Swap)
		// Elegimos dos puntos al azar
		p1 := rand.Intn(nCities-1)
		p2 := rand.Intn(nCities-p1-1) + p1 + 1

		// Invertimos el segmento para crear el hijo
		a, b := p1, p2
		for a < b {
			tourHijo[a], tourHijo[b] = tourHijo[b], tourHijo[a]
			a++
			b--
		}

		hijo := Plancton{
			Tour: tourHijo,
			Cost: utils.CalcularCostoPermutacion(tourHijo, cities),
		}
		hijos = append(hijos, hijo)
	}

	// 3. Retornar la población fusionada con los nuevos individuos
	return append(poblacion, hijos...)
}