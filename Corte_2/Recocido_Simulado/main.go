package main

import (
	"flag"
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
	"tsp-sa/parser"
	"tsp-sa/simulatedannealing"
	"tsp-sa/solver"
	"tsp-sa/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	initialTemp := flag.Float64("temp", 1000.0, "Temperatura Inicial del Recocido")
	alpha := flag.Float64("alpha", 0.995, "Factor de enfriamiento (Alpha)")
	minTemp := flag.Float64("min_temp", 0.001, "Temperatura mínima de parada")
	iterPerTemp := flag.Int("iter", 1000, "Iteraciones por nivel de temperatura")

	// Parsear los argumentos de la línea de comandos
	flag.Parse()

	// Ruta por defecto o por argumento
	archivo := "../Benchmark/berlin52.tsp" 
	args := flag.Args()
	if len(args) > 0 {
		archivo = args[0]
	}

	// Imprimir configuración para verificar
	fmt.Printf("Configuración SA: Temp=%.2f, Alpha=%.4f, Min=%.4f, Iter=%d\n", 
		*initialTemp, *alpha, *minTemp, *iterPerTemp)	

	// 1. Leer Archivo
	ciudades, err := parser.LeerArchivoTSP(archivo)
	if err != nil {
		fmt.Printf("ERROR: No se pudo leer el archivo.\n")
		fmt.Printf("Detalle: %v\n", err)
		fmt.Println("Verificar que la carpeta 'Benchmark' exista.")
		return
	}

	configSA := simulatedannealing.SAConfig{
		InitialTemp: *initialTemp,
		Alpha:       *alpha,
		MinTemp:     *minTemp,
		IterPerTemp: *iterPerTemp,
	}
	start := time.Now()

	// Ejecutar Algoritmo
	mejorTourLS, mejorCostoLS := solver.LocalSearch(ciudades)
	mejorTourSA, mejorCostoSA := solver.SimulatedAnnealingSolver(mejorTourLS, mejorCostoLS, configSA)
	
	elapsed := time.Since(start)

	// CÁLCULO DEL GAP
	optimo := utils.GetOptimalCost(archivo)
	//gapLS := 0.0
	gapSA := 0.0

	if optimo > 0 {
		//gapLS = (mejorCostoLS - optimo) / optimo * 100
		gapSA = (mejorCostoSA - optimo) / optimo * 100
	}

	// Imprimimos en formato tabla 
	nombreArchivo := filepath.Base(archivo)

	_ = mejorTourSA

	
	fmt.Printf("Benchmark: %s\n", nombreArchivo)
	fmt.Printf("%-10s\t%-10s\t%-6s\t%-10s\n", "Tiempo", "Costo", "Optimo", "GAP SA (%)")
	fmt.Printf("%s\t%.4f\t%.0f\t%.2f\n", elapsed, mejorCostoSA, optimo, gapSA)

}