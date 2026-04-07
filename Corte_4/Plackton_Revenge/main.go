package main

import (
	"flag"
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
	"tsp/parser"
	"tsp/plancton"
	"tsp/utils"
)

func main() {
	// Inicializar la semilla aleatoria
	rand.Seed(time.Now().UnixNano())

	// Configurar flags de la OFP en la terminal
	pop := flag.Int("pop", 50, "Tamaño de la poblacion (N)")
	iter := flag.Int("iter", 1000, "Numero maximo de iteraciones")
	alpha := flag.Float64("alpha", 0.1, "Intensidad de corrientes para Deriva (Alpha)")
	delta := flag.Float64("delta", 50.0, "Paso quimiotactico inicial / Vecinos a explorar (Delta)")
	gamma := flag.Float64("gamma", 0.95, "Decaimiento quimiotactico (Gamma)")
	bloom := flag.Float64("bloom", 0.1, "Porcentaje de florecimiento (BloomPct)")
	tfreq := flag.Int("tfreq", 50, "Frecuencia de turbulencia en iteraciones (T)")
	tmu := flag.Float64("tmu", 0.2, "Intensidad de turbulencia / Fraccion perturbada (Mu)")
	flat := flag.Bool("flat", false, "Mostrar informacion en formato plano (sin encabezados)")

	// Parsear los argumentos de la linea de comandos
	flag.Parse()

	// Ruta por defecto o por argumento
	archivo := "../Benchmark/berlin52.tsp"
	args := flag.Args()
	if len(args) > 0 {
		archivo = args[0]
	}

	// 1. Leer Archivo usando tu parser original
	ciudades, err := parser.LeerArchivoTSP(archivo)
	if err != nil {
		fmt.Printf("ERROR: No se pudo leer el archivo.\n")
		fmt.Printf("Detalle: %v\n", err)
		fmt.Println("Verificar que la carpeta 'Benchmark' exista.")
		return
	}

	// 2. Configurar los parámetros de la Metaheurística
	configOFP := plancton.OFPConfig{
		PopSize:    *pop,
		MaxIter:    *iter,
		Alpha:      *alpha,
		DeltaInit:  *delta,
		Gamma:      *gamma,
		BloomPct:   *bloom,
		TurbFreq:   *tfreq,
		TurbIntens: *tmu,
	}

	// 3. Ejecutar OFP y medir el tiempo
	start := time.Now()
	result := plancton.EjecutarOFP(ciudades, configOFP)
	elapsed := time.Since(start)

	// 4. Calculo del GAP con tu BKS
	optimo := utils.GetOptimalCost(archivo)
	gapOFP := 0.0
	if optimo > 0 {
		gapOFP = (result.BestCost - optimo) / optimo * 100
	}

	// 5. Imprimir resultados
	nombreArchivo := filepath.Base(archivo)

	if *flat {
		// Formato CSV para scripts (ej. run_benchmarks.sh)
		fmt.Printf("%s,%.4f,%s,%.0f,%.2f,%d,%d,%.2f,%.2f,%.2f,%.2f,%d,%.2f,%d\n",
			nombreArchivo, result.BestCost, elapsed, optimo, gapOFP,
			*pop, *iter, *alpha, *delta, *gamma, *bloom, *tfreq, *tmu,
			result.LastImproveGen)
	} else {
		// Formato legible para consola
		fmt.Printf("%-10s\t%-10s\t%-10s\t%-6s\t%-10s\n",
			"Benchmark", "Tiempo", "Costo", "Optimo", "GAP OFP (%)")
		fmt.Printf("%s\t%s\t%.4f\t%.0f\t%.2f\n",
			nombreArchivo, elapsed, result.BestCost, optimo, gapOFP)
		fmt.Printf("Config OFP: Pop=%d, Iter=%d, Alpha=%.2f, Delta=%.2f, Gamma=%.2f, Bloom=%.2f, TFreq=%d, TMu=%.2f\n",
			*pop, *iter, *alpha, *delta, *gamma, *bloom, *tfreq, *tmu)
		fmt.Printf("Convergencia: ultima mejora en iteracion %d\n", result.LastImproveGen)
	}
}
