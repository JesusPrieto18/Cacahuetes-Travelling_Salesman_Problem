package main

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"os"
	"time"
	"tsp-ils/parser"
	"tsp-ils/solver"
	"tsp-ils/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Ruta por defecto o por argumento
	archivo := "../Benchmark/berlin52.tsp" 
	if len(os.Args) > 1 {
		archivo = os.Args[1]
	}

	//fmt.Println("=============================================")
	//fmt.Printf(" PROYECTO TSP - LS Solver with 2-Opt\n")
	//fmt.Printf(" Buscando archivo: %s\n", archivo)
	//fmt.Println("=============================================")

	// 1. Leer Archivo
	ciudades, err := parser.LeerArchivoTSP(archivo)
	if err != nil {
		fmt.Printf("ERROR: No se pudo leer el archivo.\n")
		fmt.Printf("Detalle: %v\n", err)
		fmt.Println("Verificar que la carpeta 'Benchmark' exista.")
		return
	}
	//fmt.Printf("Cargado correctamente: %d ciudades.\n", len(ciudades))
	//fmt.Println("---------------------------------------------")

	start := time.Now()

	// 2. Ejecutar Algoritmo
	mejorTourLS, mejorCostoLS := solver.LocalSearch(ciudades)
	mejorTourSA, mejorCostoSA := solver.SimulatedAnnealingSolver(mejorTourLS, mejorCostoLS)
	
	elapsed := time.Since(start)

	// 3. CÁLCULO DEL GAP
	optimo := utils.GetOptimalCost(archivo)
	//gapLS := 0.0
	gapSA := 0.0

	if optimo > 0 {
		//gapLS = (mejorCostoLS - optimo) / optimo * 100
		gapSA = (mejorCostoSA - optimo) / optimo * 100
	}

	// Imprimimos en formato tabla para comparar
	// Formato: Archivo | TiempoTotal | CostoLS (Gap) | CostoSA (Gap) | Optimo
	nombreArchivo := filepath.Base(archivo)
	/*
	fmt.Println("---------------------------------------------------------------------------------")
	fmt.Printf("Instancia: %s\n", nombreArchivo)
	fmt.Printf("Tiempo Total: %s\n", elapsed)
	fmt.Println("---------------------------------------------------------------------------------")
	fmt.Printf("%-10s | %-10s | %-10s\n", "Metodo", "Costo", "GAP (%)")
	fmt.Println("---------------------------------------------------------------------------------")
	
	if optimo > 0 {
		fmt.Printf("2-Opt      | %.4f  | %.2f%%\n", mejorCostoLS, gapLS)
		fmt.Printf("Sim. Ann.  | %.4f  | %.2f%%\n", mejorCostoSA, gapSA)
		fmt.Printf("BKS        | %.0f       | 0.00%%\n", optimo)
	} else {
		fmt.Printf("2-Opt      | %.4f  | N/A\n", mejorCostoLS)
		fmt.Printf("Sim. Ann.  | %.4f  | N/A\n", mejorCostoSA)
	}
	fmt.Println("---------------------------------------------------------------------------------")
	*/
	// Evitar errores de variable no usada
	_ = mejorTourSA

	
	fmt.Printf("%s\n", nombreArchivo)
	fmt.Printf("%-10s\t%-10s\t%-6s\t%-10s\n", "Tiempo", "Costo", "Optimo", "GAP SA (%)")
	fmt.Printf("%s\t%.4f\t%.0f\t%.2f\n", elapsed, mejorCostoSA, optimo, gapSA)
	//fmt.Printf("%s\t%s\n", nombreArchivo, elapsed)
	//fmt.Printf("%s\t%.4f\n", nombreArchivo, mejorCostoSA)
	//fmt.Printf("%s\t%.0f\n", nombreArchivo, optimo)
	//fmt.Printf("%s\t%.2f%%\n", nombreArchivo, gapSA)
	
	//fmt.Println("---------------------------------------------")

}