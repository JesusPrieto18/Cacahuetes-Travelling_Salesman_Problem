package main

import (
	"fmt"
	"flag"
	"math/rand"
	"path/filepath"
	"time"
	"tsp-common/parser"
	"tsp-common/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	popSize    := flag.Int("pop", 30, "Tamaño de la población")
	maxGen     := flag.Int("gen", 1000, "Número máximo de generaciones")
	mutRate    := flag.Float64("mut", 0.15, "Probabilidad de mutación (doble-puente)")
	nParents   := flag.Int("parents", 3, "Número de padres para recombinación (≥3)")
	convThresh := flag.Int("conv", 3, "Umbral de distancia promedio para reinicio")
	flat       := flag.Bool("flat", false, "Mostrar información en formato plano (sin encabezados)")

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
		fmt.Printf("ERROR: El archivo '%s' no contiene ciudades (NODE_COORD_SECTION no encontrado o vacío).\n", archivo)
		return
	}

	dist := buildDistMatrix(cities)

	ma := &MA{
		dist:       dist,
		cities:     cities,
		popSize:    *popSize,
		maxGen:     *maxGen,
		mutRate:    *mutRate,
		nParents:   *nParents,
		convThresh: *convThresh,
	}

	start := time.Now()

	// Ejecutar algoritmo memético
	bestTour, bestCost := ma.Run()

	elapsed := time.Since(start)

	// Calcular GAP
	optimo := utils.GetOptimalCost(archivo)
	gap := 0.0
	if optimo > 0 {
		gap = (bestCost - optimo) / optimo * 100
	}

	_ = bestTour

	nombreArchivo := filepath.Base(archivo)

	// Imprimir resultados
	if *flat {
		fmt.Printf("%s\t%s\t%.4f\t%.0f\t%.2f\t%d\t%d\t%.4f\t%d\n",
			nombreArchivo, elapsed, bestCost, optimo, gap,
			*popSize, *maxGen, *mutRate, *nParents)
	} else {
		fmt.Printf("%-10s\t%-10s\t%-10s\t%-6s\t%-10s\n",
			"Benchmark", "Tiempo", "Costo", "Optimo", "GAP MA (%)")
		fmt.Printf("%s\t%s\t%.4f\t%.0f\t%.2f\n",
			nombreArchivo, elapsed, bestCost, optimo, gap)
		fmt.Printf("Configuración MA: Pop=%d, Gen=%d, Mut=%.4f, Padres=%d, Conv=%d\n",
			*popSize, *maxGen, *mutRate, *nParents, *convThresh)
	}
}