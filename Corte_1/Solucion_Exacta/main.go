package main

import (
	"math"
	"sort"
	"container/heap"

	"fmt"
	"strings"

	"flag"
	"os"
	"path/filepath"
	"time"

	"solucion_exacta/tsp"
)


// Calcula el Lower Bound (límite inferior) para un nodo dado.
// 
// Estrategia:
// 1. Para cada ciudad, sumamos las dos aristas de menor costo (entrada y salida)
// 2. Para ciudades en el camino fijo, usamos las aristas comprometidas
// 3. Dividimos entre 2 porque cada arista se cuenta dos veces
//
// Distances: una matriz de slices de flotantes con las distancias entre cada par de nodos
// path: el camino actual, cada ciudad es un entero
// visited: un diccionario donde la llave es el entero que representa la ciudad e indica si ha sido visitada o no
// currentCost: el costo actual del tour
func CalculateLowerBound(distances [][]float64, path []int, visited map[int]bool, currentCost float64) float64 {

	n := len(distances)

	// Empezamos con el costo del camino ya recorrido
	lb := currentCost

	// Para el cálculo optimista del costo restante
	minCostSum := 0.0

	// Identificar ciudades no visitadas
	unvisited := make([]int, 0)

	for i := 0; i < n; i++ {
		if !visited[i] {
			unvisited = append(unvisited, i)
		}
	}

	// Para cada ciudad no visitada, sumamos sus dos aristas más baratas
	for _, city := range unvisited {

		// Obtenemos todas las distancias desde/hacia esta ciudad
		edges := make([]float64, 0, n-1)

		for otherCity := 0; otherCity < n; otherCity++ {
			if otherCity != city {
				edges = append(edges, distances[city][otherCity])
			}
		}
		
		sort.Float64s(edges)
		// Sumamos las dos más pequeñas (estrategia de aproximación de 1-Tree reducida)
		if len(edges) >= 2 {
			minCostSum += edges[0] + edges[1]
		}
	}

	// Lógica para conectar el camino actual con el futuro y el regreso al inicio

	// Para la última ciudad en el path actual, necesitamos considerar 
    // una arista de salida hacia ciudades no visitadas
	if len(path) > 0 && len(unvisited) > 0 {

		lastCity := path[len(path)-1]
		
		// Mínimo desde la última ciudad del path a una no visitada
		minToUnvisited := math.MaxFloat64
		for _, v := range unvisited {
			if distances[lastCity][v] < minToUnvisited {
				minToUnvisited = distances[lastCity][v]
			}
		}

		minCostSum += minToUnvisited

		// Mínimo desde una no visitada de regreso al inicio (ciudad 0)
		minToStart := math.MaxFloat64
		for _, v := range unvisited {
			if distances[v][0] < minToStart {
				minToStart = distances[v][0]
			}
		}
		minCostSum += minToStart
	}

	// Cada arista se cuenta dos veces (ida y vuelta conceptual)
	return lb + (minCostSum / 2.0)

}


// Item representa un nodo en el árbol de búsqueda 
// Lo usaremos en TSPBranchBoundWithLB
type Item struct {
	lowerBound  float64
	currentCity int
	path        []int
	visited     map[int]bool
	actualCost  float64

	// Requerido por la interfaz heap
	index       int 
}

// PriorityQueue implementa heap.Interface y contiene Items
// Cola de prioridad que usaremos en TSPBranchBoundWithLB
type PriorityQueue []*Item

// Funciones que requiere la cola de prioridad
func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].lowerBound < pq[j].lowerBound }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}


// Resuelve el TSP usando Branch and Bound con Lower Bound y Best-First Search.
    
// Args:
//     distances: Matriz de adyacencia con los pesos de las aristas
//     node_names: Lista con los nombres de los nodos (opcional)
// Retorna:
//     - best_path: Lista con el orden óptimo de nodos
//     - best_cost: Costo total del tour óptimo    
func TSPBranchBoundWithLB(distances [][]float64) ([]int, float64) {

	n := len(distances)
	var bestPath []int
	bestCost := math.Inf(1)

	// <----- Aqui va una logica de strings que voy a quitar --->

	// Nodo inicial visitado: el 0
	initialVisited := map[int]bool{0: true}
	initialPath := []int{0}

	// Calculamos el LB inicial para el nodo raíz
	initialLB := CalculateLowerBound(distances, initialPath, initialVisited, 0)


	// Usamos una cola de prioridad (min-heap) para Best-First Search
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{
		lowerBound:  initialLB,
		currentCity: 0,
		path:        initialPath,
		visited:     initialVisited,
		actualCost:  0,
	})

	nodesExplored := 0
	nodesPruned := 0

	println("Iniciando busqueda")
	var iter = 0

	for pq.Len() > 0 {
		
		fmt.Printf("Iteracion: %d\n", iter)
		iter++
		fmt.Printf("#Items en la cola de prioridad: %d\n", pq.Len())

		// Extraemos el nodo con menor Lower Bound (Best-First Search)
		// Hacer pop en nuestra implementacion devuelve interface{}, hay que convertirlo a *Item
		node := heap.Pop(pq).(*Item)
		nodesExplored++

		// Poda: Si el LB de este nodo es mayor o igual al mejor costo conocido,
        // no puede contener una solución mejor

		if node.lowerBound >= bestCost {
			nodesPruned++
			continue
		}

		// Si hemos visitado todas las ciudades, completamos el tour
		if len(node.path) == n {
			// Agregamos el costo de regresar al inicio
			totalCost := node.actualCost + distances[node.currentCity][0]

			// Actualizamos la mejor solución si encontramos una mejor
			if totalCost < bestCost {
				bestPath = append([]int(nil), node.path...) // Copia profunda del slice
				bestCost = totalCost
				
				fmt.Printf("Nueva mejor solución costo: %.2f | Nodos: %d\n", bestCost, nodesExplored)
			}
		} else 
		{
			// Explorar hijos
			// Generamos nodos hijos para todas las ciudades no visitadas

			for nextCity := 0; nextCity < n; nextCity++ {
				if node.visited[nextCity]{
					continue;
				}
				
				// Calculamos el costo acumulado al ir a la siguiente ciudad
				newActualCost := node.actualCost + distances[node.currentCity][nextCity]

				// Poda temprana: si el costo acumulado ya supera el mejor costo, no expandimos
				if newActualCost >= bestCost {
					nodesPruned++
					continue
				}

				// Crear nuevo path y visited para el hijo
				newPath := append([]int(nil), node.path...) // Se podria mejorar 
				newPath = append(newPath, nextCity)
				
				newVisited := make(map[int]bool)
				for k, v := range node.visited {
					newVisited[k] = v
				}

				newVisited[nextCity] = true

				// Calculamos el Lower Bound para este nodo hijo
				newLB := CalculateLowerBound(distances, newPath, newVisited, newActualCost)

				// Poda: solo agregamos el nodo si su LB es mejor que el mejor costo actual
				if newLB < bestCost {
					heap.Push(pq, &Item{
						lowerBound:  newLB,
						currentCity: nextCity,
						path:        newPath,
						visited:     newVisited,
						actualCost:  newActualCost,
					})
				} else {
					nodesPruned++
				}
				
			}
		}


	}
	
	fmt.Printf("\nFinalizado. Nodos explorados: %d, Podados: %d\n", nodesExplored, nodesPruned)
	fmt.Printf("\nTotal de nodos considerados: %d\n", nodesExplored + nodesPruned)


	return bestPath, bestCost

}

func main() {
	// CLI flags
	tspFile := flag.String("tsp", "", "Path to TSPLIB .tsp file (e.g., berlin52.tsp)")

	flag.Parse()

	if *tspFile == "" {
		fmt.Fprintf(os.Stderr, "Error: must specify -tsp <file.tsp>\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -tsp <file.tsp>\n", os.Args[0])
		os.Exit(1)
	}

	// Load TSPLIB instance
	inst, err := tsp.LoadTSPLIB(*tspFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading instance: %v\n", err)
		os.Exit(1)
	}

	instanceName := strings.TrimSuffix(filepath.Base(*tspFile), ".tsp")
	

	fmt.Printf("\nInstance: %s (%d cities)\n", instanceName, inst.NumCities)
	if inst.OptimalCost > 0 {
		fmt.Printf("Optimal cost: %.0f\n", inst.OptimalCost)
	} else {
		fmt.Println("Optimal cost: unknown")
	}
	fmt.Println()

	// Run 
	start := time.Now()

	bestPath, bestCost := TSPBranchBoundWithLB(inst.Distance)

	elapsed := time.Since(start)


	// Calculate gap
	gap := 0.0
	if inst.OptimalCost > 0 {
		gap = (bestCost - inst.OptimalCost) / inst.OptimalCost * 100
	}

	// Print results
	fmt.Println("Results:")
	fmt.Printf("  Best tour length: %.0f\n", bestCost)
	if inst.OptimalCost > 0 {
		fmt.Printf("  Optimal length: %.0f\n", inst.OptimalCost)
		fmt.Printf("  Gap: %.2f%%\n", gap)
	}

	fmt.Printf("  Time: %v\n", elapsed)
	fmt.Printf("  Tour: %v\n", bestPath)
	println("")
}