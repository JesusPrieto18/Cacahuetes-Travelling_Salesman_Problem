package plancton

// OFPConfig almacena los parámetros entonables de la Metaheurística
// Optimización por Florecimiento de Plancton (OFP).
type OFPConfig struct {
	PopSize int // N: Cantidad de planctones (Tamaño de población)
	MaxIter int // Condición de parada (Equivalente a número máximo de iteraciones)

	// Operador 1 - Deriva
	Alpha float64 // α: Intensidad de corrientes (0.0 a 1.0 para determinar el tamaño del salto)

	// Operador 2 - Quimiotaxis
	DeltaInit float64 // δ inicial: Paso quimiotáctico (Tamaño inicial del vecindario en búsqueda local)
	Gamma     float64 // Γ: Tasa de decaimiento quimiotáctico (0.0 a 1.0). Para la fórmula: δ_t = δ_0 * (Γ^t)

	// Operador 3 - Florecimiento (Bloom)
	BloomPct float64 // B: Porcentaje de florecimiento (Fracción élite que se reproduce, ej. 0.1 para 10%)

	// Operador 4 - Turbulencia
	TurbFreq   int     // T: Frecuencia de turbulencia (Cada cuántas iteraciones ocurre la dispersión)
	TurbIntens float64 // μ: Intensidad de turbulencia (Fracción de la población que será perturbada)
}

// Plancton representa a un individuo (solución candidata) en el entorno.
type Plancton struct {
	Tour []int   // Representación discreta: permutación de índices de ciudades
	Cost float64 // Fitness: costo total de la ruta (a minimizar)
}
