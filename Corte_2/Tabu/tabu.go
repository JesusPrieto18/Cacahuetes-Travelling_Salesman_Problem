package main

import (
	"math"
	"math/rand"
	"tsp-common/models"
	"tsp-common/utils"
)

func TabuSearch(ciudades []models.City, maxIteraciones int, tenenciaTabu int) ([]models.City, float64) {
	n := len(ciudades)

	// 1. Solución Inicial (Aleatoria o Greedy)
	tourActual := utils.CopiarTour(ciudades)
	rand.Shuffle(len(tourActual), func(i, j int) {
		tourActual[i], tourActual[j] = tourActual[j], tourActual[i]
	})

	costoActual := utils.CalcularCostoTotal(tourActual)

	// Mejor solución global (Best Global)
	tourBest := utils.CopiarTour(tourActual)
	costoBest := costoActual

	// 2. Estructura de Memoria Tabú
	maxID := 0
	for _, c := range ciudades {
		if c.ID > maxID {
			maxID = c.ID
		}
	}
	tabuMatrix := make([][]int, maxID+1)
	for i := range tabuMatrix {
		tabuMatrix[i] = make([]int, maxID+1)
	}

	// 3. Bucle Principal
	for iter := 1; iter <= maxIteraciones; iter++ {

		// Variables para encontrar el MEJOR vecino en TODA la vecindad (Best Improvement)
		mejorVecinoCosto := math.MaxFloat64
		moveI, moveJ := -1, -1

		foundMove := false

		// Explorar toda la vecindad 2-Opt
		for i := 1; i < n-1; i++ {
			for j := i + 1; j < n; j++ {

				// A. Calcular Delta (Diferencia de costo)
				d1 := utils.DistanciaEuclidiana(tourActual[i-1], tourActual[i])
				d2 := utils.DistanciaEuclidiana(tourActual[j], tourActual[(j+1)%n])
				costoAristasViejas := d1 + d2

				// Costo nuevo de las aristas a crear
				d3 := utils.DistanciaEuclidiana(tourActual[i-1], tourActual[j])
				d4 := utils.DistanciaEuclidiana(tourActual[i], tourActual[(j+1)%n])
				costoAristasNuevas := d3 + d4

				delta := costoAristasNuevas - costoAristasViejas
				nuevoCostoPosible := costoActual + delta

				// B. Verificar Estado Tabú
				id1 := tourActual[i].ID
				id2 := tourActual[j].ID

				esTabu := false
				if tabuMatrix[id1][id2] > iter {
					esTabu = true
				}

				// C. Criterio de Aspiración
				if esTabu && (nuevoCostoPosible < costoBest) {
					esTabu = false // Romper el tabú por aspiración
				}

				// D. Selección del mejor vecino admisible
				if !esTabu {
					if nuevoCostoPosible < mejorVecinoCosto {
						mejorVecinoCosto = nuevoCostoPosible
						moveI = i
						moveJ = j
						foundMove = true
					}
				}
			}
		}

		// 4. Moverse a la siguiente solución
		if foundMove {
			invertirSegmento(tourActual, moveI, moveJ)
			costoActual = mejorVecinoCosto

			// Actualizar la lista Tabú
			id1 := tourActual[moveI].ID
			id2 := tourActual[moveJ].ID

			tabuMatrix[id1][id2] = iter + tenenciaTabu
			tabuMatrix[id2][id1] = iter + tenenciaTabu

			// Actualizar el mejor global si corresponde
			if costoActual < costoBest {
				tourBest = utils.CopiarTour(tourActual)
				costoBest = costoActual
			}
		}
	}

	return tourBest, costoBest
}
