package geneticalgorithm

import (
	"tsp-ds/models"
)

// PathRelinking genera un camino de soluciones desde el tourInicial hasta el tourGuia
// introduciendo gradualmente las aristas del tourGuia. Devuelve el mejor tour intermedio.
func PathRelinking(tourInicial, tourGuia []int, cities []models.City) []int {
	n := len(tourInicial)

	// 1. Guardar el mejor encontrado en el trayecto
	mejorIntermedio := make([]int, n)
	copy(mejorIntermedio, tourInicial)
	mejorCosto := EvaluateCost(mejorIntermedio, cities)

	// 2. Tour actual que se irá transformando paso a paso
	actual := make([]int, n)
	copy(actual, tourInicial)

	// 3. Extraer aristas del guía para búsqueda rápida (O(1))
	aristasGuia := make(map[[2]int]bool)
	for i := 0; i < n; i++ {
		u := tourGuia[i]
		v := tourGuia[(i+1)%n]
		if u > v { u, v = v, u }
		aristasGuia[[2]int{u, v}] = true
	}

	// 4. Transformación gradual (máximo N pasos)
	for pasos := 0; pasos < n; pasos++ {
		var uTarget, vTarget int
		faltaArista := false

		// a) Buscar una arista de la Guía que le falte al Actual
		for edge := range aristasGuia {
			tieneArista := false
			for i := 0; i < n; i++ {
				uAct := actual[i]
				vAct := actual[(i+1)%n]
				if uAct > vAct { uAct, vAct = vAct, uAct }
				
				if edge[0] == uAct && edge[1] == vAct {
					tieneArista = true
					break
				}
			}
			
			if !tieneArista {
				uTarget = edge[0]
				vTarget = edge[1]
				faltaArista = true
				break // Tomamos la primera arista faltante que encontremos
			}
		}

		// Si ya no faltan aristas, hemos llegado exactamente al tourGuia
		if !faltaArista {
			break
		}

		// b) Insertar la arista (uTarget, vTarget) en 'actual' forzando un 2-opt
		idxU, idxV := -1, -1
		for i := 0; i < n; i++ {
			if actual[i] == uTarget { idxU = i }
			if actual[i] == vTarget { idxV = i }
		}

		// Rotamos el arreglo para que uTarget esté en la posición 0 (facilita la inversión)
		actual = rotarCero(actual, idxU)
		
		// Buscamos dónde quedó vTarget después de rotar
		for i := 0; i < n; i++ {
			if actual[i] == vTarget { 
				idxV = i
				break 
			}
		}

		// Invertimos el segmento desde el índice 1 hasta idxV.
		// ¡Esto conecta automáticamente actual[0] (uTarget) con actual[1] (vTarget)!
		invertirSegmentoRelink(actual, 1, idxV)

		// c) Evaluar la nueva solución intermedia generada
		costoActual := EvaluateCost(actual, cities)
		if costoActual < mejorCosto {
			mejorCosto = costoActual
			copy(mejorIntermedio, actual)
		}
	}

	return mejorIntermedio
}

// rotarCero desplaza los elementos del arreglo circularmente para que el elemento
// en 'nuevaRaiz' quede en la posición 0.
func rotarCero(tour []int, nuevaRaiz int) []int {
	n := len(tour)
	rotado := make([]int, n)
	for i := 0; i < n; i++ {
		rotado[i] = tour[(nuevaRaiz+i)%n]
	}
	return rotado
}

// invertirSegmentoRelink invierte un subarreglo en su lugar.
func invertirSegmentoRelink(tour []int, i, j int) {
	for i < j {
		tour[i], tour[j] = tour[j], tour[i]
		i++
		j--
	}
}