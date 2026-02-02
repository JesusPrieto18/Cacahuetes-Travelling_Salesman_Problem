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

	fmt.Println("=============================================")
	fmt.Printf(" PROYECTO TSP - LS Solver with 2-Opt\n")
	fmt.Printf(" Buscando archivo: %s\n", archivo)
	fmt.Println("=============================================")

	// 1. Leer Archivo
	ciudades, err := parser.LeerArchivoTSP(archivo)
	if err != nil {
		fmt.Printf("ERROR: No se pudo leer el archivo.\n")
		fmt.Printf("Detalle: %v\n", err)
		fmt.Println("Verificar que la carpeta 'Benchmark' exista.")
		return
	}
	fmt.Printf("Cargado correctamente: %d ciudades.\n", len(ciudades))
	fmt.Println("---------------------------------------------")

	start := time.Now()

	// 2. Ejecutar Algoritmo
	mejorTour, mejorCosto := solver.LocalSearch(ciudades)

	elapsed := time.Since(start)

	// 3. CÁLCULO DEL GAP
	optimo := utils.GetOptimalCost(archivo)
	gap := 0.0
	
	if optimo > 0 {
		gap = (mejorCosto - optimo) / optimo * 100
	}

	// 4. Reporte de Resultados
	fmt.Println("---------------------------------------------")
	fmt.Printf("Tiempo Total:    %s\n", elapsed)
	fmt.Printf("Mejor Costo:     %.4f\n", mejorCosto)
	
	if optimo > 0 {
		fmt.Printf("Óptimo (BKS):    %.0f\n", optimo)
		fmt.Printf("GAP:             %.2f%%\n", gap)
		
		// Interpretación rápida
		if gap < 0.01 {
			fmt.Println(">> ¡Resultado Óptimo encontrado!")
		} else if gap < 5.0 {
			fmt.Println(">> Resultado de alta calidad.")
		}
	} else {
		fmt.Println("GAP: Desconocido (Instancia no registrada)")
	}
	fmt.Println("---------------------------------------------")
	// Evitar error de variable no usada si no imprimes el tour
	_ = mejorTour

	nombreArchivo := filepath.Base(archivo)

	fmt.Printf("%s\t%s\n", nombreArchivo, elapsed)
	fmt.Printf("%s\t%.4f\n", nombreArchivo, mejorCosto)
	fmt.Printf("%s\t%.0f\n", nombreArchivo, optimo)
	fmt.Printf("%s\t%.2f%%\n", nombreArchivo, gap)
	
	fmt.Println("---------------------------------------------")

}