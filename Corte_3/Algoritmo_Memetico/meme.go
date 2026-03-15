package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"tsp-common/models"
	"tsp-common/utils"
	"tsp-sa/localsearch"
)

// Tour es una permutación de índices 0..n-1 sobre el slice de ciudades.
type Tour []int

func (t Tour) clone() Tour {
	c := make(Tour, len(t))
	copy(c, t)
	return c
}

// Matriz de distancias 
type DistMatrix [][]float64

func buildDistMatrix(cities []models.City) DistMatrix {
	n := len(cities)
	d := make(DistMatrix, n)
	for i := range d {
		d[i] = make([]float64, n)
		for j := range d[i] {
			d[i][j] = utils.DistanciaEuclidiana(cities[i], cities[j])
		}
	}
	return d
}

func tourCost(dist DistMatrix, t Tour) float64 {
	n := len(t)
	c := 0.0
	for i := 0; i < n; i++ {
		c += dist[t[i]][t[(i+1)%n]]
	}
	return c
}

// tourToCities convierte un Tour (índices) al slice de Cities que usa utils.
func tourToCities(t Tour, cities []models.City) []models.City {
	out := make([]models.City, len(t))
	for i, idx := range t {
		out[i] = cities[idx]
	}
	return out
}

// Distancia entre tours
// d(t1,t2) = número de aristas en t1 que NO están en t2.
func tourDist(a, b Tour) int {
	n := len(a)
	edgesB := make(map[[2]int]bool, n)
	for i := 0; i < n; i++ {
		u, v := b[i], b[(i+1)%n]
		if u > v {
			u, v = v, u
		}
		edgesB[[2]int{u, v}] = true
	}
	diff := 0
	for i := 0; i < n; i++ {
		u, v := a[i], a[(i+1)%n]
		if u > v {
			u, v = v, u
		}
		if !edgesB[[2]int{u, v}] {
			diff++
		}
	}
	return diff
}

// Búsqueda local: 2-opt (wrapper sobre localsearch.TwoOpt)
// TwoOpt trabaja con []models.City, así que convertimos Tour ↔ []City.
func applyTwoOpt(cities []models.City, t Tour) (Tour, float64) {
	cityTour := tourToCities(t, cities)
	improved, cost := localsearch.TwoOpt(cityTour)
	return citiesToTour(improved, cities), cost
}

// citiesToTour convierte []City de vuelta a Tour (índices) usando el ID de cada ciudad.
func citiesToTour(cityTour []models.City, cities []models.City) Tour {
	// índice inverso: ID → posición en el slice original
	idToIdx := make(map[int]int, len(cities))
	for i, c := range cities {
		idToIdx[c.ID] = i
	}
	t := make(Tour, len(cityTour))
	for i, c := range cityTour {
		t[i] = idToIdx[c.ID]
	}
	return t
}

// Recombinación respetuosa
// Implementa la recombinación descripta en la Clase 10:
//   - RESPETUOSA: las aristas comunes a TODOS los padres aparecen en el hijo.
//   - Los fragmentos se unen con greedy (vecino más cercano),
//     prefiriendo aristas que no estén en ningún padre (DPX-like).
// Soporta n ≥ 2 padres: las aristas comunes son la intersección de todos.

func commonEdgesAll(parents []Tour) map[[2]int]bool {
	if len(parents) == 0 {
		return nil
	}
	n := len(parents[0])

	// Aristas del primer padre
	common := make(map[[2]int]bool, n)
	for i := 0; i < n; i++ {
		u, v := parents[0][i], parents[0][(i+1)%n]
		if u > v {
			u, v = v, u
		}
		common[[2]int{u, v}] = true
	}

	// Intersectar con cada padre restante
	for _, par := range parents[1:] {
		edgesPar := make(map[[2]int]bool, n)
		for i := 0; i < n; i++ {
			u, v := par[i], par[(i+1)%n]
			if u > v {
				u, v = v, u
			}
			edgesPar[[2]int{u, v}] = true
		}
		for e := range common {
			if !edgesPar[e] {
				delete(common, e)
			}
		}
	}
	return common
}

func recombine(dist DistMatrix, parents []Tour) Tour {
	n := len(parents[0])
	common := commonEdgesAll(parents)

	// Grafo de adyacencia con solo las aristas comunes
	adj := make([][]int, n)
	deg := make([]int, n)
	for e := range common {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
		deg[e[0]]++
		deg[e[1]]++
	}

	// Extraer fragmentos (cadenas de aristas comunes)
	visited := make([]bool, n)
	var fragments [][]int

	for start := 0; start < n; start++ {
		if visited[start] || deg[start] != 1 {
			continue
		}
		frag := []int{start}
		visited[start] = true
		cur := start
		for {
			moved := false
			for _, nb := range adj[cur] {
				if !visited[nb] {
					frag = append(frag, nb)
					visited[nb] = true
					cur = nb
					moved = true
					break
				}
			}
			if !moved {
				break
			}
		}
		fragments = append(fragments, frag)
	}
	// Nodos aislados (grado 0) → fragmentos unitarios
	for i := 0; i < n; i++ {
		if !visited[i] {
			fragments = append(fragments, []int{i})
			visited[i] = true
		}
	}

	// Caso degenerado: un único fragmento ya es el tour completo
	if len(fragments) == 1 && len(fragments[0]) == n {
		return Tour(fragments[0])
	}

	// Aristas prohibidas = unión de aristas de todos los padres
	// (evitar re-insertar aristas que diferencian a los padres)
	parentEdges := make(map[[2]int]bool)
	for _, par := range parents {
		for i := 0; i < n; i++ {
			u, v := par[i], par[(i+1)%n]
			if u > v {
				u, v = v, u
			}
			parentEdges[[2]int{u, v}] = true
		}
	}

	// Conectar fragmentos con greedy (estilo DPX)
	for len(fragments) > 1 {
		frag0 := fragments[0]
		endpoint := frag0[len(frag0)-1]

		bestCost := math.MaxFloat64
		bestNode := -1
		bestFrag := -1

		for fi := 1; fi < len(fragments); fi++ {
			f := fragments[fi]
			for _, cand := range []int{f[0], f[len(f)-1]} {
				e := [2]int{min2(endpoint, cand), max2(endpoint, cand)}
				d := dist[endpoint][cand]
				if parentEdges[e] {
					d += 1e9 // penalizar aristas ya en algún padre
				}
				if d < bestCost {
					bestCost = d
					bestNode = cand
					bestFrag = fi
				}
			}
		}

		target := fragments[bestFrag]
		if target[len(target)-1] == bestNode {
			for l, r := 0, len(target)-1; l < r; l, r = l+1, r-1 {
				target[l], target[r] = target[r], target[l]
			}
		}
		merged := append(frag0, target...)
		fragments[0] = merged
		fragments = append(fragments[:bestFrag], fragments[bestFrag+1:]...)
	}

	result := Tour(fragments[0])
	// Completar si faltan nodos (caso raro)
	if len(result) != n {
		inResult := make([]bool, n)
		for _, v := range result {
			inResult[v] = true
		}
		for i := 0; i < n; i++ {
			if !inResult[i] {
				result = append(result, i)
			}
		}
	}
	return result
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max2(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Mutación: doble-puente (4-cambio no secuencial)
func doubleBridge(t Tour) Tour {
	n := len(t)
	pts := rand.Perm(n)[:4]
	sort.Ints(pts)
	a, b, c := pts[0], pts[1], pts[2]
	result := make(Tour, 0, n)
	result = append(result, t[:a+1]...)
	result = append(result, t[b+1:c+1]...)
	result = append(result, t[a+1:b+1]...)
	result = append(result, t[c+1:]...)
	return result
}

// Individuo 
type Individual struct {
	tour Tour
	cost float64
}

// Inicialización
func initPopulation(dist DistMatrix, cities []models.City, size int) []Individual {
	n := len(dist)
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}
	seen := make(map[string]bool)
	pop := make([]Individual, 0, size)
	attempts := 0
	for len(pop) < size && attempts < size*20 {
		attempts++
		rand.Shuffle(n, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
		t := Tour(append([]int{}, perm...))
		t, _ = applyTwoOpt(cities, t)
		key := fmt.Sprint(t)
		if !seen[key] {
			seen[key] = true
			pop = append(pop, Individual{tour: t, cost: tourCost(dist, t)})
		}
	}
	sort.Slice(pop, func(i, j int) bool { return pop[i].cost < pop[j].cost })
	return pop
}

// Algoritmo Memético
type MA struct {
	dist       DistMatrix
	cities     []models.City
	popSize    int
	maxGen     int
	mutRate    float64
	nParents   int // ≥ 3 (requerimiento del proyecto)
	convThresh int // distancia promedio mínima antes de reiniciar
}

func (ma *MA) Run() (Tour, float64) {
	pop := initPopulation(ma.dist, ma.cities, ma.popSize)
	best := pop[0]

	for gen := 0; gen < ma.maxGen; gen++ {
		// Selección de nParents padres distintos al azar
		idx := rand.Perm(len(pop))[:ma.nParents]
		parents := make([]Tour, ma.nParents)
		for i, pi := range idx {
			parents[i] = pop[pi].tour.clone()
		}

		// Recombinación respetuosa de ≥3 padres
		child := recombine(ma.dist, parents)

		// Mutación (doble-puente) para diversificación
		if rand.Float64() < ma.mutRate {
			child = doubleBridge(child)
		}

		// Mejora local post-recombinación → intensificación (núcleo del AM)
		var childCost float64
		child, childCost = applyTwoOpt(ma.cities, child)

		// Reemplazo elitista: entra si mejora al peor
		worst := len(pop) - 1
		if childCost < pop[worst].cost {
			pop[worst] = Individual{tour: child, cost: childCost}
			sort.Slice(pop, func(i, j int) bool { return pop[i].cost < pop[j].cost })
		}

		if pop[0].cost < best.cost {
			best = pop[0]
			fmt.Printf("  Gen %d → mejor costo: %.2f\n", gen, best.cost)
		}

		// Reinicio si la población convergió
		if ma.converged(pop) {
			pop = ma.restart(pop)
		}
	}
	return best.tour, best.cost
}

func (ma *MA) converged(pop []Individual) bool {
	total, pairs := 0, 0
	for i := 0; i < len(pop)-1; i++ {
		for j := i + 1; j < len(pop); j++ {
			total += tourDist(pop[i].tour, pop[j].tour)
			pairs++
		}
	}
	if pairs == 0 {
		return false
	}
	return total/pairs < ma.convThresh
}

func (ma *MA) restart(pop []Individual) []Individual {
	newPop := []Individual{pop[0]} // conservar el mejor
	for i := 1; i < len(pop); i++ {
		t := doubleBridge(pop[i].tour.clone())
		t, _ = applyTwoOpt(ma.cities, t)
		newPop = append(newPop, Individual{tour: t, cost: tourCost(ma.dist, t)})
	}
	sort.Slice(newPop, func(i, j int) bool { return newPop[i].cost < newPop[j].cost })
	return newPop
}