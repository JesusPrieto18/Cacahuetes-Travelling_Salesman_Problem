package plancton

import (
	"sort"
	"tsp/models"
	"tsp/utils"
)

// OFPResult estructura los resultados finales de la metaheurística.
type OFPResult struct {
	BestTour       []models.City // La ruta final decodificada en ciudades
	BestCost       float64       // El mejor fitness encontrado
	LastImproveGen int           // Iteración de la última mejora
	TotalIter      int           // Iteraciones ejecutadas
}

// EjecutarOFP orquesta el ciclo de vida de la Optimización por Florecimiento de Plancton.
func EjecutarOFP(oceano utils.Oceano, config OFPConfig) OFPResult {
	nCities := len(oceano)
	
	// 1. Inicialización
	poblacion := InicializarPoblacion(oceano, config.PopSize)
	
	// Rastrear el mejor global
	mejorGlobal := Plancton{
		Tour: utils.CopiarPermutacion(poblacion[0].Tour),
		Cost: poblacion[0].Cost,
	}
	lastImprove := 0
	
	// Paso quimiotáctico inicial (que irá decayendo)
	deltaActual := config.DeltaInit

	// Bucle Generacional (El paso del tiempo en el océano)
	for t := 0; t < config.MaxIter; t++ {
		
		// OPERADOR 1: Deriva (Corrientes arrastran al plancton)
		// Empezamos en i = 1 para proteger a la Élite (índice 0) de la destrucción
		for i := 1; i < len(poblacion); i++ {
			AplicarDeriva(&poblacion[i], config.Alpha)
		}

		// OPERADOR 2: Quimiotaxis (Búsqueda local de nutrientes)
		for i := range poblacion {
			AplicarQuimiotaxis(&poblacion[i], oceano, deltaActual)
		}

		// OPERADOR 3: Florecimiento / Bloom (Intensificación)
		poblacion = AplicarFlorecimiento(poblacion, config.BloomPct, oceano)

		// OPERADOR 4: Turbulencia (Diversificación periódica)
		if t > 0 && t%config.TurbFreq == 0 {
			AplicarTurbulencia(poblacion, config.TurbIntens, oceano)
		}

		// OPERADOR 5: Hundimiento (Selección natural)
		// Ordenar de mejor a peor fitness (menor costo primero)
		sort.Slice(poblacion, func(i, j int) bool {
			return poblacion[i].Cost < poblacion[j].Cost
		})
		
		// El plancton que no encuentra nutrientes y queda fuera de N, muere y se hunde.
		// En Go esto se hace eficientemente truncando el slice.
		poblacion = poblacion[:config.PopSize]

		// Actualizar el mejor global encontrado
		if poblacion[0].Cost < mejorGlobal.Cost {
			mejorGlobal.Cost = poblacion[0].Cost
			mejorGlobal.Tour = utils.CopiarPermutacion(poblacion[0].Tour)
			lastImprove = t + 1
		}

		// Enfriamiento del paso quimiotáctico (simula movimientos más finos)
		deltaActual *= config.Gamma
	}

	// Traducir el genotipo (permutación de índices) al fenotipo (slice de ciudades)
	bestTourCities := make([]models.City, nCities)
	for i, idx := range mejorGlobal.Tour {
		bestTourCities[i] = oceano[idx]
	}

	// Devolver el resultado formateado
	return OFPResult{
		BestTour:       bestTourCities,
		BestCost:       mejorGlobal.Cost,
		LastImproveGen: lastImprove,
		TotalIter:      config.MaxIter,
	}
}