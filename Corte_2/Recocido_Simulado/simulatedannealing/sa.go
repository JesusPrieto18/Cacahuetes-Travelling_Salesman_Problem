package simulatedannealing

import (
	"math"
	"math/rand"
	"tsp-ils/models"
	"tsp-ils/utils"
)

// Configuración de parámetros para SA
type SAConfig struct {
	InitialTemp float64 // Temperatura inicial
	Alpha       float64 // Factor de enfriamiento (ej. 0.99)
	MinTemp     float64 // Temperatura de parada
	IterPerTemp int     // Iteraciones por cada nivel de temperatura (Equilibrio térmico)
}

// EjecutarSA aplica Recocido Simulado sobre un tour existente
func EjecutarSA(tourInicial []models.City, config SAConfig) ([]models.City, float64) {
	
	// Inicialización
	tourActual := utils.CopiarTour(tourInicial)
	costoActual := utils.CalcularCostoTotal(tourActual)

	mejorTour := utils.CopiarTour(tourActual)
	mejorCosto := costoActual

	tempActual := config.InitialTemp
	n := len(tourActual)

	// 2. Bucle principal de temperatura
	for tempActual > config.MinTemp {
		
		// 3. Equilibrio térmico (Iteraciones a temperatura constante)
		for k := 0; k < config.IterPerTemp; k++ {
			
			// A. Generar vecino aleatorio (Movimiento 2-Opt aleatorio)
			// Seleccionamos dos índices al azar i y j
			i := rand.Intn(n)
			j := rand.Intn(n)

			// Aseguramos que i < j para facilitar la inversión
			if i > j {
				i, j = j, i
			}
			// Evitamos invertir el tour completo o segmentos nulos
			if i == j || (i == 0 && j == n-1) {
				continue
			}

			// B. Calcular Delta E (Cambio de costo)
			// Optimización: Calculamos solo la diferencia de aristas, no el tour completo
			// Aristas eliminadas: (i-1 -> i) y (j -> j+1)
			// Aristas nuevas:     (i-1 -> j) y (i -> j+1)
			
			idxPrevI := (i - 1 + n) % n
			idxNextJ := (j + 1) % n

			dEliminada1 := utils.DistanciaEuclidiana(tourActual[idxPrevI], tourActual[i])
			dEliminada2 := utils.DistanciaEuclidiana(tourActual[j], tourActual[idxNextJ])
			
			dNueva1 := utils.DistanciaEuclidiana(tourActual[idxPrevI], tourActual[j])
			dNueva2 := utils.DistanciaEuclidiana(tourActual[i], tourActual[idxNextJ])

			delta := (dNueva1 + dNueva2) - (dEliminada1 + dEliminada2)

			// Criterio de Aceptación (Metropolis)
			aceptar := false
			if delta < 0 {
				// Mejora: Aceptar siempre
				aceptar = true
			} else {
				// Empeora: Aceptar con probabilidad e^(-delta/T)
				prob := math.Exp(-delta / tempActual)
				if rand.Float64() < prob {
					aceptar = true
				}
			}

			// Aplicar cambio si es aceptado
			if aceptar {
				invertirSegmento(tourActual, i, j)
				costoActual += delta

				// Actualizar el mejor global encontrado
				if costoActual < mejorCosto {
					mejorCosto = costoActual
					mejorTour = utils.CopiarTour(tourActual)
				}
			}
		}

		// 4. Enfriamiento
		tempActual *= config.Alpha
	}

	return mejorTour, mejorCosto
}

// Función auxiliar para invertir segmento (igual que en 2-opt, pero local)
func invertirSegmento(tour []models.City, i, j int) {
	for i < j {
		tour[i], tour[j] = tour[j], tour[i]
		i++
		j--
	}
}