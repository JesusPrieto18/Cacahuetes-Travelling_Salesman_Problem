package grasp

import (
	"math/rand"
	"sort"
	"tsp-sa/models"
	"tsp-sa/utils"
)

type Candidate struct {
	city  models.City
	dist  float64
	index int
}

type Alpha struct { //
	Value   float64
	costSum float64
	uses    int
	prob    float64
}

func buildRCL(last models.City, unvisited []models.City, alpha float64) []Candidate {
	if len(unvisited) == 0 {
		return nil
	}

	candidates := make([]Candidate, len(unvisited))
	minDist, maxDist := 1e18, -1e18

	for i, c := range unvisited {
		d := utils.DistanciaEuclidiana(last, c)
		candidates[i] = Candidate{c, d, i}
		if d < minDist {
			minDist = d
		}
		if d > maxDist {
			maxDist = d
		}
	}

	// Ordenar candidatos por distancia para facilitar el sesgo
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].dist < candidates[j].dist
	})

	threshold := minDist + alpha*(maxDist-minDist)
	var rcl []Candidate
	for _, cand := range candidates {
		if cand.dist <= threshold {
			rcl = append(rcl, cand)
		} else {
			break
		}
	}
	return rcl
}

func buildSolution(cities []models.City, alpha float64) []models.City {
	n := len(cities)
	startIdx := rand.Intn(n)
	tour := []models.City{cities[startIdx]}

	// Preparar lista de no visitadas excluyendo la inicial
	unvisited := make([]models.City, 0, n-1)
	for i, c := range cities {
		if i != startIdx {
			unvisited = append(unvisited, c)
		}
	}

	for len(unvisited) > 0 {
		last := tour[len(tour)-1]

		rcl := buildRCL(last, unvisited, alpha)

		selected := chooseWithBias(rcl)

		tour = append(tour, selected.city)

		// Remover la ciudad seleccionada de la lista de no visitadas
		for i, c := range unvisited {
			if c.ID == selected.city.ID {
				unvisited = append(unvisited[:i], unvisited[i+1:]...)
				break
			}
		}
	}
	return tour
}

func chooseWithBias(rcl []Candidate) Candidate {
	// La RCL ya viene ordenada por distancia
	n := len(rcl)
	if n == 1 {
		return rcl[0]
	}

	sumWeights := n * (n + 1) / 2
	r := rand.Intn(sumWeights) + 1

	currentSum := 0
	for i := range n {
		currentSum += (n - i)
		if r <= currentSum {
			return rcl[i]
		}
	}
	return rcl[0]
}
