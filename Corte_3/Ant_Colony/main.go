package main

import (
	"flag"
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
	"tsp-common/parser"
	"tsp-common/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	numAnts := flag.Int("ants", 30, "Número de hormigas")
	numIter := flag.Int("gen", 1000, "Número de iteraciones")
	alpha := flag.Float64("alpha", 1.0, "Parámetro que pesa el nivel de feromona")
	beta := flag.Float64("beta", 5.0, "Parámetro que pesa la información heurística (1/d)")
	evap := flag.Float64("evap", 0.5, "Tasa de evaporación de feromona (rho)")
	q := flag.Float64("q", 100.0, "Constante para el depósito de feromona (Q)")
	flat := flag.Bool("flat", false, "Mostrar información en formato plano (sin encabezados)")

	flag.Parse()

	archivo := "../Benchmark/berlin52.tsp"
	args := flag.Args()
	if len(args) > 0 {
		archivo = args[0]
	}

	// Leer archivo
	cities, err := parser.LeerArchivoTSP(archivo)
	if err != nil {
		fmt.Printf("ERROR: No se pudo leer el archivo.\n")
		fmt.Printf("Detalle: %v\n", err)
		fmt.Println("Verificar que la carpeta 'Benchmark' exista.")
		return
	}
	if len(cities) == 0 {
		fmt.Printf("ERROR: El archivo '%s' no contiene ciudades.\n", archivo)
		return
	}

	aco := NewACO(cities, *numAnts, *numIter, *alpha, *beta, *evap, *q)

	start := time.Now()

	bestTour, bestCost := aco.Run()

	elapsed := time.Since(start)

	optimo := utils.GetOptimalCost(archivo)
	gap := 0.0
	if optimo > 0 {
		gap = (bestCost - optimo) / optimo * 100
	}

	_ = bestTour

	nombreArchivo := filepath.Base(archivo)

	if *flat {
		fmt.Printf("%s\t%s\t%.4f\t%.0f\t%.2f\t%d\t%d\t%.2f\t%.2f\t%.2f\t%.2f\n",
			nombreArchivo, elapsed, bestCost, optimo, gap,
			*numAnts, *numIter, *alpha, *beta, *evap, *q)
	} else {
		fmt.Printf("%-10s\t%-10s\t%-10s\t%-6s\t%-10s\n",
			"Benchmark", "Tiempo", "Costo", "Optimo", "GAP ACO (%)")
		fmt.Printf("%s\t%s\t%.4f\t%.0f\t%.2f\n",
			nombreArchivo, elapsed, bestCost, optimo, gap)
		fmt.Printf("Configuración ACO: Hormigas=%d, Gen=%d, Alpha=%.2f, Beta=%.2f, Evap=%.2f, Q=%.2f\n",
			*numAnts, *numIter, *alpha, *beta, *evap, *q)
	}
}
