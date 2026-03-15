package main

import (
	"fmt"
	"math"
	"math/rand"
	"tsp-common/models"
	"tsp-common/utils"
)

type ACO struct {
	dist        [][]float64
	cities      []models.City
	numAnts     int
	numIter     int
	alpha       float64
	beta        float64
	evaporation float64
	q           float64
	pheromone   [][]float64
	heuristic   [][]float64
}

type Ant struct {
	path    []int
	visited []bool
	cost    float64
}

func buildDistMatrix(cities []models.City) [][]float64 {
	n := len(cities)
	d := make([][]float64, n)
	for i := range d {
		d[i] = make([]float64, n)
		for j := range d[i] {
			d[i][j] = utils.DistanciaEuclidiana(cities[i], cities[j])
		}
	}
	return d
}

func NewACO(cities []models.City, numAnts, numIter int, alpha, beta, evaporation, q float64) *ACO {
	n := len(cities)
	dist := buildDistMatrix(cities)

	pheromone := make([][]float64, n)
	heuristic := make([][]float64, n)

	// Inicializar feromonas y matriz heurística
	initialPheromone := 1e-6
	for i := 0; i < n; i++ {
		pheromone[i] = make([]float64, n)
		heuristic[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			pheromone[i][j] = initialPheromone
			if i != j {
				if dist[i][j] > 0 {
					heuristic[i][j] = 1.0 / dist[i][j]
				} else {
					heuristic[i][j] = 1.0 / 1e-8
				}
			}
		}
	}

	return &ACO{
		dist:        dist,
		cities:      cities,
		numAnts:     numAnts,
		numIter:     numIter,
		alpha:       alpha,
		beta:        beta,
		evaporation: evaporation,
		q:           q,
		pheromone:   pheromone,
		heuristic:   heuristic,
	}
}

func (aco *ACO) Run() ([]int, float64) {
	n := len(aco.cities)
	bestCost := math.MaxFloat64
	var bestPath []int

	for iter := 0; iter < aco.numIter; iter++ {
		ants := make([]Ant, aco.numAnts)
		for k := 0; k < aco.numAnts; k++ {
			ants[k] = aco.buildAntSolution()
			if ants[k].cost < bestCost {
				bestCost = ants[k].cost
				bestPath = make([]int, n)
				copy(bestPath, ants[k].path)
				fmt.Printf("  Iter %d → Hormiga %d mejor costo encontrado: %.2f\n", iter, k, bestCost)
			}
		}
		aco.updatePheromones(ants)
	}

	return bestPath, bestCost
}

func (aco *ACO) buildAntSolution() Ant {
	n := len(aco.cities)
	ant := Ant{
		path:    make([]int, 0, n),
		visited: make([]bool, n),
	}

	// Comienza en una ciudad al azar
	start := rand.Intn(n)
	ant.path = append(ant.path, start)
	ant.visited[start] = true

	curr := start
	for step := 1; step < n; step++ {
		next := aco.selectNextCity(curr, ant.visited)
		ant.path = append(ant.path, next)
		ant.visited[next] = true
		curr = next
	}

	// Calcular costo del recorrido completo
	cost := 0.0
	for i := 0; i < n-1; i++ {
		cost += aco.dist[ant.path[i]][ant.path[i+1]]
	}
	cost += aco.dist[ant.path[n-1]][ant.path[0]]
	ant.cost = cost

	return ant
}

func (aco *ACO) selectNextCity(curr int, visited []bool) int {
	n := len(aco.cities)
	probs := make([]float64, n)
	sumProbs := 0.0

	for i := 0; i < n; i++ {
		if !visited[i] {
			p := math.Pow(aco.pheromone[curr][i], aco.alpha) * math.Pow(aco.heuristic[curr][i], aco.beta)
			probs[i] = p
			sumProbs += p
		}
	}

	// Fallback si la suma de probabilidades es muy cercana a 0 debido al underflow
	if sumProbs <= 0.0 {
		var unvisited []int
		for i := 0; i < n; i++ {
			if !visited[i] {
				unvisited = append(unvisited, i)
			}
		}
		return unvisited[rand.Intn(len(unvisited))]
	}

	// Elegir ruta basados en las probabilidades 
	r := rand.Float64() * sumProbs
	acc := 0.0
	for i := 0; i < n; i++ {
		if !visited[i] {
			acc += probs[i]
			if r <= acc {
				return i
			}
		}
	}

	// Fallback para problemas de precisión
	for i := 0; i < n; i++ {
		if !visited[i] {
			return i
		}
	}
	return -1
}

func (aco *ACO) updatePheromones(ants []Ant) {
	n := len(aco.cities)

	// Fase de Evaporación
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			aco.pheromone[i][j] *= (1.0 - aco.evaporation)
			// Evitar que la feromona caiga estrictamente a cero
			if aco.pheromone[i][j] < 1e-12 {
				aco.pheromone[i][j] = 1e-12
			}
		}
	}

	// Fase de Depósito de feromona
	for _, ant := range ants {
		deposit := aco.q / ant.cost
		for i := 0; i < n-1; i++ {
			from, to := ant.path[i], ant.path[i+1]
			aco.pheromone[from][to] += deposit
			aco.pheromone[to][from] += deposit  // Suponiendo problema de ruta simétrica
		}
		// Regreso a la base
		from, to := ant.path[n-1], ant.path[0]
		aco.pheromone[from][to] += deposit
		aco.pheromone[to][from] += deposit
	}
}
