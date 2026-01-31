package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"tsp-ils/parser"
	"tsp-ils/solver"
)

// Funcion main
func main() {
	rand.Seed(time.Now().UnixNano())

	rutaPorDefecto := "../Benchmark/berlin52.tsp"
	archivo := rutaPorDefecto

	// Si pasas un argumento por consola, usa ese en su lugar
	if len(os.Args) > 1 {
		archivo = os.Args[1]
	}

	fmt.Println("=============================================")
	fmt.Printf(" PROYECTO TSP - ILS Solver\n")
	fmt.Printf(" Buscando archivo en: %s\n", archivo)
	fmt.Println("=============================================")

	// 1. Leer Archivo
	ciudades, err := parser.LeerArchivoTSP(archivo)
	if err != nil {
		fmt.Printf("ERROR CRÍTICO: No se pudo leer el archivo.\n")
		fmt.Printf("Detalle: %v\n", err)
		fmt.Println("Verificar que la carpeta 'Benchmark' exista.")
		return
	}
	fmt.Printf("Cargado correctamente: %d ciudades.\n", len(ciudades))
	fmt.Println("---------------------------------------------")

	start := time.Now()

	// 2. Ejecutar Algoritmo
	mejorTour, mejorCosto := solver.ILS(ciudades, 3000)

	elapsed := time.Since(start)

	// 3. Resultados
	fmt.Println("---------------------------------------------")
	fmt.Printf("Tiempo de ejecución: %s\n", elapsed)
	fmt.Printf("MEJOR COSTO FINAL: %.4f\n", mejorCosto)
	fmt.Println("---------------------------------------------")

	fmt.Print("Ruta Final: ")
	for i, c := range mejorTour {
		if i < 15 {
			fmt.Printf("%d -> ", c.ID)
		}
	}
	fmt.Println("... ->", mejorTour[0].ID)
}