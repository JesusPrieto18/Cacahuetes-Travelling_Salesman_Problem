package main

import (
	"flag"
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
	"tsp-ga/geneticalgorithm"
	"tsp-ga/parser"
	"tsp-ga/solver"
	"tsp-ga/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	pop := flag.Int("pop", 600, "Tamaño de la poblacion")
	gen := flag.Int("gen", 2000, "Numero maximo de generaciones")
	mut := flag.Float64("mut", 0.3, "Probabilidad de mutacion")
	tourn := flag.Int("tourn", 3, "Tamaño del torneo para seleccion")
	stag := flag.Int("stag", 200, "Generaciones sin mejora antes de parar (0 = desactivado)")
	flat := flag.Bool("flat", false, "Mostrar informacion en formato plano (sin encabezados)")

	// Parsear los argumentos de la linea de comandos
	flag.Parse()

	// Ruta por defecto o por argumento
	archivo := "../Benchmark/berlin52.tsp"
	args := flag.Args()
	if len(args) > 0 {
		archivo = args[0]
	}

	// 1. Leer Archivo
	ciudades, err := parser.LeerArchivoTSP(archivo)
	if err != nil {
		fmt.Printf("ERROR: No se pudo leer el archivo.\n")
		fmt.Printf("Detalle: %v\n", err)
		fmt.Println("Verificar que la carpeta 'Benchmark' exista.")
		return
	}

	configGA := geneticalgorithm.GAConfig{
		PopSize:         *pop,
		Generations:     *gen,
		MutationRate:    *mut,
		TournamentSize:  *tourn,
		StagnationLimit: *stag,
	}

	start := time.Now()

	// 2. Ejecutar Algoritmo Genetico
	result := solver.GeneticAlgorithmSolver(ciudades, configGA)

	elapsed := time.Since(start)

	// 3. Calculo del GAP
	optimo := utils.GetOptimalCost(archivo)
	gapGA := 0.0
	if optimo > 0 {
		gapGA = (result.BestCost - optimo) / optimo * 100
	}

	// 4. Imprimir resultados
	nombreArchivo := filepath.Base(archivo)

	if *flat {
		fmt.Printf("%s,%.4f,%s,%.0f,%.2f,%d,%d,%.4f,%d,%d,%d,%d,%s\n",
			nombreArchivo, result.BestCost, elapsed, optimo, gapGA,
			*pop, *gen, *mut, *tourn, *stag,
			result.LastImproveGen, result.TotalGens, result.StopReason)
	} else {
		fmt.Printf("%-10s\t%-10s\t%-10s\t%-6s\t%-10s\n",
			"Benchmark", "Tiempo", "Costo", "Optimo", "GAP GA (%)")
		fmt.Printf("%s\t%s\t%.4f\t%.0f\t%.2f\n",
			nombreArchivo, elapsed, result.BestCost, optimo, gapGA)
		fmt.Printf("Configuracion GA: Pop=%d, Gen=%d, Mut=%.4f, Tourn=%d, Stag=%d\n",
			*pop, *gen, *mut, *tourn, *stag)
		fmt.Printf("Convergencia: ultima mejora en gen %d, parada en gen %d por %s\n",
			result.LastImproveGen, result.TotalGens, result.StopReason)
	}
}
