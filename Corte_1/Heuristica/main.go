package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"heuristica/acs"
	"heuristica/tsp"
)

func main() {
	defaults := acs.DefaultParams()

	// CLI flags
	tspFile := flag.String("tsp", "", "Path to TSPLIB .tsp file (e.g., berlin52.tsp)")
	iterations := flag.Int("iterations", defaults.Iterations, "Number of ACS iterations")
	numAnts := flag.Int("ants", defaults.NumAnts, "Number of ants")
	beta := flag.Float64("beta", defaults.Beta, "Beta parameter (heuristic weight)")
	rho := flag.Float64("rho", defaults.Rho, "Rho parameter (local evaporation)")
	alpha := flag.Float64("alpha", defaults.Alpha, "Alpha parameter (global evaporation)")
	q0 := flag.Float64("q0", defaults.Q0, "Q0 parameter (exploitation probability)")
	seed := flag.Int64("seed", 0, "Random seed (0 = use current time)")
	verbose := flag.Bool("verbose", false, "Print detailed output")

	flag.Parse()

	if *tspFile == "" {
		fmt.Fprintf(os.Stderr, "Error: must specify -tsp <file.tsp>\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -tsp <file.tsp> [options]\n", os.Args[0])
		os.Exit(1)
	}

	// Load TSPLIB instance
	inst, err := tsp.LoadTSPLIB(*tspFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading instance: %v\n", err)
		os.Exit(1)
	}

	instanceName := strings.TrimSuffix(filepath.Base(*tspFile), ".tsp")

	if *verbose {
		fmt.Printf("Instance: %s (%d cities)\n", instanceName, inst.NumCities)
		if inst.OptimalCost > 0 {
			fmt.Printf("Optimal cost: %.0f\n", inst.OptimalCost)
		} else {
			fmt.Println("Optimal cost: unknown")
		}
		fmt.Println()
	}

	// Set up parameters
	params := acs.Params{
		NumAnts:    *numAnts,
		Beta:       *beta,
		Rho:        *rho,
		Alpha:      *alpha,
		Q0:         *q0,
		Iterations: *iterations,
	}

	// Set seed
	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}

	if *verbose {
		fmt.Println("ACS Parameters:")
		fmt.Printf("  Ants: %d\n", params.NumAnts)
		fmt.Printf("  Beta: %.2f\n", params.Beta)
		fmt.Printf("  Rho: %.2f\n", params.Rho)
		fmt.Printf("  Alpha: %.2f\n", params.Alpha)
		fmt.Printf("  Q0: %.2f\n", params.Q0)
		fmt.Printf("  Iterations: %d\n", params.Iterations)
		fmt.Printf("  Seed: %d\n", *seed)
		fmt.Println()
	}

	// Run ACS
	start := time.Now()
	solver := acs.New(inst, params, *seed)
	bestTour, bestLength := solver.Run()
	elapsed := time.Since(start)

	// Calculate gap
	gap := 0.0
	if inst.OptimalCost > 0 {
		gap = (bestLength - inst.OptimalCost) / inst.OptimalCost * 100
	}

	// Output results
	if *verbose {
		fmt.Println("Results:")
		fmt.Printf("  Best tour length: %.0f\n", bestLength)
		if inst.OptimalCost > 0 {
			fmt.Printf("  Optimal length: %.0f\n", inst.OptimalCost)
			fmt.Printf("  Gap: %.2f%%\n", gap)
		}
		fmt.Printf("  Time: %v\n", elapsed)
		fmt.Printf("  Tour: %v\n", bestTour)
	} else {
		// Compact output for scripting
		if inst.OptimalCost > 0 {
			fmt.Printf("Instance: %s, Cities: %d, Best: %.0f, Optimal: %.0f, Gap: %.2f%%, Time: %v\n",
				instanceName, inst.NumCities, bestLength, inst.OptimalCost, gap, elapsed)
		} else {
			fmt.Printf("Instance: %s, Cities: %d, Best: %.0f, Time: %v\n",
				instanceName, inst.NumCities, bestLength, elapsed)
		}
	}
}
