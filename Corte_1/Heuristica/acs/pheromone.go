package acs

// PheromoneMatrix manages the pheromone levels between cities
type PheromoneMatrix struct {
	tau  [][]float64 // pheromone levels
	tau0 float64     // initial pheromone value
}

// NewPheromoneMatrix creates a new pheromone matrix initialized with tau0
// Where Lnn is the nearest neighbor tour length
func NewPheromoneMatrix(n int, lnn float64) *PheromoneMatrix {
	tau0 := 1.0 / (float64(n) * lnn)

	tau := make([][]float64, n)
	for i := range tau {
		tau[i] = make([]float64, n)
		for j := range tau[i] {
			tau[i][j] = tau0
		}
	}

	return &PheromoneMatrix{
		tau:  tau,
		tau0: tau0,
	}
}

// Get returns the pheromone level between cities i and j
func (pm *PheromoneMatrix) Get(i, j int) float64 {
	return pm.tau[i][j]
}

// LocalUpdate applies the local pheromone update rule after an ant moves
func (pm *PheromoneMatrix) LocalUpdate(i, j int, rho float64) {
	pm.tau[i][j] = (1-rho)*pm.tau[i][j] + rho*pm.tau0
	pm.tau[j][i] = pm.tau[i][j] // symmetric
}

// GlobalUpdate applies the global pheromone update at the end of an iteration
// Only the best ant deposits pheromone
func (pm *PheromoneMatrix) GlobalUpdate(tour []int, tourLength float64, alpha float64) {
	deposit := alpha / tourLength
	n := len(tour)

	for i := 0; i < n; i++ {
		from := tour[i]
		to := tour[(i+1)%n]
		pm.tau[from][to] = (1-alpha)*pm.tau[from][to] + deposit
		pm.tau[to][from] = pm.tau[from][to] // symmetric
	}
}

// Tau0 returns the initial pheromone value
func (pm *PheromoneMatrix) Tau0() float64 {
	return pm.tau0
}
