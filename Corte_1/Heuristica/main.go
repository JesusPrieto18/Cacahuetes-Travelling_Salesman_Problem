package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"heuristica/tsp"
)

func main() {
	// CLI flags
	tspFile := flag.String("tsp", "", "Path to TSPLIB .tsp file (e.g., berlin52.tsp)")
	verbose := flag.Bool("verbose", false, "Print detailed output")

	flag.Parse()

	if *tspFile == "" {
		fmt.Fprintf(os.Stderr, "Error: must specify -tsp <file.tsp>\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -tsp <file.tsp> [-verbose]\n", os.Args[0])
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
		fmt.Printf("Algorithm: Farthest Insertion\n")
		fmt.Println()
	}

	// Run Farthest Insertion heuristic
	start := time.Now()
	bestTour, bestLength := tsp.FarthestInsertion(inst)
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
