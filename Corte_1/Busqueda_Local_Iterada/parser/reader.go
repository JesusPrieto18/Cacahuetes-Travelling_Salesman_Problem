package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"tsp-ils/models"
)

// Funcion para leer el archivo TSP
func LeerArchivoTSP(rutaArchivo string) ([]models.City, error) {
	// Intentamos abrir el archivo en la ruta especificada
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

		// Detectar inicio de sección de coordenadas
		if line == "NODE_COORD_SECTION" {
			readingCoords = true
			continue
		}

		// Leer datos: ID X Y
		if readingCoords {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				id, err1 := strconv.Atoi(fields[0])
				x, err2 := strconv.ParseFloat(fields[1], 64)
				y, err3 := strconv.ParseFloat(fields[2], 64)

				if err1 == nil && err2 == nil && err3 == nil {
					cities = append(cities, models.City{ID: id, X: x, Y: y})
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cities, nil
}
