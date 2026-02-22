package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"
	"tsp-sa/grasp"
	"tsp-sa/parser"
	"tsp-sa/utils"
)

func main() {
	// Configuracion inicial y semilla de aleatoriedad
	rand.Seed(time.Now().UnixNano())
	file := "../Benchmark/bier127.tsp"
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		file = args[0]
	}

	// Leer archivo
	cities, err := parser.LeerArchivoTSP(file)
	if err != nil {
		fmt.Printf("ERROR: No se pudo leer el archivo.\n")
		fmt.Printf("Detalle: %v\n", err)
		fmt.Println("Verificar que la carpeta 'Benchmark' exista.")
		return
	}

	// Inicializar alpha
	var alpha grasp.Alpha = grasp.Alpha{Value: 0.0} // si no se pasa argumento, se usara 0.0 (greedy puro)
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%f", &alpha.Value)
		if alpha.Value < 0 || alpha.Value > 1 {
			fmt.Printf("Valor de alpha fuera de rango (0.0 - 1.0). Usando valor por defecto 0.0.\n")
			alpha.Value = 0.0
		}
	}

	fmt.Printf("Iniciando GRASP Reactivo para %d ciudades con alpha = %.2f...\n", len(cities), alpha.Value)

	// GraspReactivo coordinara la construccion, el sesgo, el inicio aleatorio y el 2-opt
	start := time.Now()
	bestTour, bestCost := grasp.GraspReactivo(cities, 1000, &alpha)
	elapsed := time.Since(start)

	// CALCULO DEL GAP
	optimo := utils.GetOptimalCost(file)
	gap := 0.0

	if optimo > 0 {
		gap = (bestCost - optimo) / optimo * 100
	}

	_ = bestTour

	printTable(file, elapsed, bestCost, optimo, gap)
}

func printTable(name string, tiempo time.Duration, result float64, optimo float64, gap float64) {
	w := tabwriter.NewWriter(os.Stdout, 8, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "Instancia\tTiempo\tResultado\tOptimo\tGAP (%)")

	line := fmt.Sprintf("%s\t%v\t%.2f\t%.2f\t%.2f",
		name,
		tiempo.Round(time.Millisecond), // Redondear tiempo para limpieza (quitar de ser necesario)
		result, optimo, gap)

	fmt.Fprintln(w, line)

	w.Flush()
}
