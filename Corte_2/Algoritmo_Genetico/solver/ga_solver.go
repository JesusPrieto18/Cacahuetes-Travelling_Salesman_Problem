package solver

import (
	"tsp-ga/geneticalgorithm"
	"tsp-ga/models"
)

// GeneticAlgorithmSolver executes the genetic algorithm on the given cities.
func GeneticAlgorithmSolver(ciudades []models.City, config geneticalgorithm.GAConfig) geneticalgorithm.GAResult {
	return geneticalgorithm.RunGA(ciudades, config)
}
