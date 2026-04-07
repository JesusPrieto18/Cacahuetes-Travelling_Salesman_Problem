package plancton

import (
	"math/rand"
	"tsp/models"
	"tsp/utils"
)

// AplicarQuimiotaxis realiza una búsqueda local acotada (Explotación).
// Evalúa 'delta' vecinos usando movimientos 2-opt y se mueve si hay mejora.
func AplicarQuimiotaxis(p *Plancton, cities []models.City, delta float64) {
	n := len(p.Tour)
	numVecinos := int(delta)

	if numVecinos < 1 || n < 4 {
		// Actualizamos el costo usando la función de tu paquete utils
		p.Cost = utils.CalcularCostoPermutacion(p.Tour, cities)
		return
	}

	// Usamos la utilidad pública para copiar
	mejorTour := utils.CopiarPermutacion(p.Tour)
	huboMejora := false

	for v := 0; v < numVecinos; v++ {
		// Elegir dos puntos de corte para generar un vecino (movimiento 2-opt)
		i := rand.Intn(n-2) + 1
		j := rand.Intn(n-i-1) + i + 1

		// Evaluar las aristas actuales vs las nuevas usando tu utils.DistanciaEuclidiana
		d1 := utils.DistanciaEuclidiana(cities[mejorTour[i-1]], cities[mejorTour[i]])
		d2 := utils.DistanciaEuclidiana(cities[mejorTour[j]], cities[mejorTour[(j+1)%n]])
		costoActualAristas := d1 + d2

		d3 := utils.DistanciaEuclidiana(cities[mejorTour[i-1]], cities[mejorTour[j]])
		d4 := utils.DistanciaEuclidiana(cities[mejorTour[i]], cities[mejorTour[(j+1)%n]])
		costoNuevoAristas := d3 + d4

		// Si mejora, aplicamos el movimiento inmediatamente en nuestra copia
		if (costoActualAristas - costoNuevoAristas) > 0.0001 {
			invertirSegmentoLocal(mejorTour, i, j)
			huboMejora = true
		}
	}

	// Sincronización final del Plancton
	if huboMejora {
		p.Tour = mejorTour
		p.Cost = utils.CalcularCostoPermutacion(mejorTour, cities)
	} else {
		p.Cost = utils.CalcularCostoPermutacion(p.Tour, cities)
	}
}

// invertirSegmentoLocal invierte un subarreglo en su lugar.
// Se mantiene privada (en minúscula) para uso exclusivo de este archivo.
func invertirSegmentoLocal(tour []int, i, j int) {
	for i < j {
		tour[i], tour[j] = tour[j], tour[i]
		i++
		j--
	}
}
