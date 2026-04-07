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
		// Usamos la utilidad pública para copiar la ruta de forma segura
		tourHijo := utils.CopiarPermutacion(padre.Tour)

		// Mutación mínima: un único swap aleatorio para buscar en la vecindad inmediata
		idx1 := rand.Intn(nCities)
		idx2 := rand.Intn(nCities)
		tourHijo[idx1], tourHijo[idx2] = tourHijo[idx2], tourHijo[idx1]

		// Crear el nuevo plancton y evaluar su costo con tu utilidad centralizada
		hijo := Plancton{
			Tour: tourHijo,
			Cost: utils.CalcularCostoPermutacion(tourHijo, cities),
		}
		hijos = append(hijos, hijo)
	}

	// 3. Retornar la población fusionada con los nuevos individuos
	return append(poblacion, hijos...)
}