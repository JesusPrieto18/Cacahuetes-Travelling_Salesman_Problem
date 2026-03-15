package geneticalgorithm

import (
	"math/rand"
	"math"
	"tsp-meme/models"
	"tsp-meme/utils"
)

// CutAndFillCrossover implements the "corte y llenado" (Order Crossover) operator.
// Selects a random cut point p. Child1 takes parent1[0..p], then fills with
// elements from parent2 in order, skipping those already present.
// Child2 is symmetric: prefix from parent2, fill from parent1.
func CutAndFillCrossover(parent1, parent2 []int) ([]int, []int) {
	n := len(parent1)

	// Cut point: at least 1 element from prefix, at least 1 to fill
	p := rand.Intn(n-1) + 1 // p in [1, n-1)

	child1 := cutAndFillBuild(parent1, parent2, p, n)
	child2 := cutAndFillBuild(parent2, parent1, p, n)

	return child1, child2
}

// cutAndFillBuild builds one child: copies prefix[0..p) from donor,
// then appends elements from filler in order, skipping duplicates.
func cutAndFillBuild(donor, filler []int, p, n int) []int {
	child := make([]int, 0, n)
	inChild := make([]bool, n)

	// Copy prefix from donor
	for i := 0; i < p; i++ {
		child = append(child, donor[i])
		inChild[donor[i]] = true
	}

	// Fill remaining from filler in order
	for _, val := range filler {
		if !inChild[val] {
			child = append(child, val)
			inChild[val] = true
		}
	}

	return child
}

// DPXMultiParentCrossover toma N padres y genera un hijo preservando las aristas comunes.
func DPXMultiParentCrossover(parents [][]int, cities []models.City) []int {
	if len(parents) == 0 {
		return nil
	}
	n := len(parents[0])

	// 1. Mapa de aristas del Padre 0 (para comparar contra el resto)
	baseEdges := make(map[[2]int]bool)
	for i := 0; i < n; i++ {
		u, v := parents[0][i], parents[0][(i+1)%n]
		if u > v { u, v = v, u }
		baseEdges[[2]int{u, v}] = true
	}

	// 2. Filtrar dejando solo las aristas que están en TODOS los padres (Consenso Estricto)
	for _, p := range parents[1:] {
		currentEdges := make(map[[2]int]bool)
		for i := 0; i < n; i++ {
			u, v := p[i], p[(i+1)%n]
			if u > v { u, v = v, u }
			currentEdges[[2]int{u, v}] = true
		}
		for edge := range baseEdges {
			if !currentEdges[edge] {
				delete(baseEdges, edge)
			}
		}
	}

	// 3. Mapa de adyacencia para el hijo basado en las aristas comunes
	adj := make([][]int, n)
	for edge := range baseEdges {
		u, v := edge[0], edge[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	// 4. Mapa de TODAS las aristas de todos los padres (para penalizarlas en la reconexión DPX)
	allParentEdges := make(map[[2]int]bool)
	for _, p := range parents {
		for i := 0; i < n; i++ {
			u, v := p[i], p[(i+1)%n]
			if u > v { u, v = v, u }
			allParentEdges[[2]int{u, v}] = true
		}
	}

	// 5. Construir el hijo (Greedy Nearest Neighbor respetando los fragmentos)
	child := make([]int, 0, n)
	visited := make([]bool, n)
	
	// Empezar desde la ciudad 0
	curr := 0
	child = append(child, curr)
	visited[curr] = true

	for len(child) < n {
		next := -1
		// Si el nodo actual tiene una conexión obligatoria (arista común), la seguimos
		for _, neighbor := range adj[curr] {
			if !visited[neighbor] {
				next = neighbor
				break
			}
		}

		// Si no hay conexión obligatoria, buscamos la ciudad no visitada más "cercana"
		if next == -1 {
			minDist := math.MaxFloat64
			for c := 0; c < n; c++ {
				if !visited[c] {
					dist := utils.DistanciaEuclidiana(cities[curr], cities[c])
					// Truco DPX de tu compañera: penalizar aristas que pertenecían a los padres
					u, v := curr, c
					if u > v { u, v = v, u }
					if allParentEdges[[2]int{u, v}] {
						dist += 1e9 
					}
					if dist < minDist {
						minDist = dist
						next = c
					}
				}
			}
		}

		curr = next
		child = append(child, curr)
		visited[curr] = true
	}

	return child
}