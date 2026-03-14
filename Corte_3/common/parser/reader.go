package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"tsp-common/models"
)

func LeerArchivoTSP(rutaArchivo string) ([]models.City, error) {
	file, err := os.Open(rutaArchivo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cities []models.City
	scanner := bufio.NewScanner(file)
	readingCoords := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "EOF" {
			break
		}
		if line == "NODE_COORD_SECTION" {
			readingCoords = true
			continue
		}

		if readingCoords {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				id, _ := strconv.Atoi(fields[0])
				x, _ := strconv.ParseFloat(fields[1], 64)
				y, _ := strconv.ParseFloat(fields[2], 64)
				cities = append(cities, models.City{ID: id, X: x, Y: y})
			}
		}
	}
	return cities, nil
}
