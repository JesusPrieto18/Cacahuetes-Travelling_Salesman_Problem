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

	maxIter := flag.Int("iter", 2000, "Máximo de iteraciones")
	tenencia := flag.Int("tenure", 25, "Tenencia Tabú")
	flat := flag.Bool("flat", false, "Mostrar informacion en formato plano (sin encabezados)")

	// Parsear los argumentos de la línea de comandos
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
		return
	}

	start := time.Now()

	// Ejecutar Algoritmo
	mejorTour, mejorCosto := TabuSearch(ciudades, *maxIter, *tenencia)

	elapsed := time.Since(start)

	// CÁLCULO DEL GAP
	optimo := utils.GetOptimalCost(archivo)
	gapTabu := 0.0

	if optimo > 0 {
		gapTabu = (mejorCosto - optimo) / optimo * 100
	}

	// Imprimimos en formato tabla
	nombreArchivo := filepath.Base(archivo)

	_ = mejorTour

	if *flat {
		fmt.Printf("%s\t%s\t%.4f\t%.0f\t%.2f\t%d\t%d\n", nombreArchivo, elapsed, mejorCosto, optimo, gapTabu, *maxIter, *tenencia)
	} else {
		fmt.Printf("%-10s\t%-10s\t%-10s\t%-6s\t%-10s\n", "Benchmark", "Tiempo", "Costo", "Optimo", "GAP Tabu (%)")
		fmt.Printf("%s\t%s\t%.4f\t%.0f\t%.2f\n", nombreArchivo, elapsed, mejorCosto, optimo, gapTabu)
		fmt.Printf("Configuración Tabu: Iter=%d, Tenure=%d\n", *maxIter, *tenencia)
	}
}
