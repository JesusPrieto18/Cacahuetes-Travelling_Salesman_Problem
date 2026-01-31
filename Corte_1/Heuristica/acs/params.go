package acs

// Params holds the parameters for the ACS algorithm
type Params struct {
	NumAnts     int     // number of ants
	Beta        float64 // weight of heuristic information
	Rho         float64 // local pheromone evaporation rate
	Alpha       float64 // global pheromone evaporation rate
	Q0          float64 // exploitation probability
	Iterations  int     // number of iterations
}

// DefaultParams returns the canonical ACS parameters
func DefaultParams() Params {
	return Params{
		NumAnts:    40,
		Beta:       2.0,
		Rho:        0.1,
		Alpha:      0.1,
		Q0:         0.9,
		Iterations: 1000,
	}
}
