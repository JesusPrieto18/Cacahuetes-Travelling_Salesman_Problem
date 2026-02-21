package main

import (
	"tsp-common/models"
)

func invertirSegmento(tour []models.City, i, j int) {
	for i < j {
		tour[i], tour[j] = tour[j], tour[i]
		i++
		j--
	}
}
