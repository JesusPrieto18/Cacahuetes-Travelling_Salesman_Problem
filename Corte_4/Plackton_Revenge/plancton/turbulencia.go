package plancton

import (
	"math/rand"
	"tsp/models"
	"tsp/utils"
)

// AplicarTurbulencia aplica una macro-mutación (Double Bridge) a una fracción de la población.
// Esto permite escapar de óptimos locales sin crear rutas "basura" que mueran instantáneamente.
func AplicarTurbulencia(poblacion []Plancton, mu float64, cities []models.City) {
	nPop := len(poblacion)
	turbCount := int(mu * float64(nPop))
	nCities := len(cities)

	// Si la población es muy pequeña o la intensidad es 0, omitimos
	if turbCount == 0 || nPop < 2 || nCities < 8 {
		return
	}

	for i := 0; i < turbCount; i++ {
		// Protegemos al campeón (índice 0)
		idx := rand.Intn(nPop-1) + 1

		// Tomamos como base la estructura del líder para no empezar desde cero (evitar muertes por costo alto)
		baseTour := poblacion[0].Tour

		// Generamos 3 puntos de corte aleatorios para dividir la ruta en 4 segmentos
		p1 := 1 + rand.Intn(nCities/4)
		p2 := p1 + 1 + rand.Intn(nCities/4)
		p3 := p2 + 1 + rand.Intn(nCities/4)

		// Ensamblamos el Doble Puente: Segmentos [A, B, C, D] se convierten en [A, D, C, B]
		nuevoTour := make([]int, 0, nCities)
		nuevoTour = append(nuevoTour, baseTour[:p1]...)       // A
		nuevoTour = append(nuevoTour, baseTour[p3:]...)       // D
		nuevoTour = append(nuevoTour, baseTour[p2:p3]...)     // C
		nuevoTour = append(nuevoTour, baseTour[p1:p2]...)     // B

		// Actualizamos el plancton arrastrado por la turbulencia
		poblacion[idx].Tour = nuevoTour
		poblacion[idx].Cost = utils.CalcularCostoPermutacion(nuevoTour, cities)
	}
}