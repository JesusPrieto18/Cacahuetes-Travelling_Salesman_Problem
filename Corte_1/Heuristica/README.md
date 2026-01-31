# ACS para TSP

Implementación de Ant Colony System (ACS) para resolver el Problema del Viajante (TSP).

## Requisitos

- Go 1.21+

## Compilar

```bash
go build -o heuristica
```

## Uso

```bash
./heuristica -tsp ../Benchmark/berlin52.tsp -verbose
```

## Flags

| Flag | Default | Descripción |
|------|---------|-------------|
| `-tsp` | - | Ruta al archivo TSPLIB (.tsp) |
| `-iterations` | `100` | Iteraciones ACS |
| `-ants` | `10` | Número de hormigas |
| `-beta` | `2.0` | Peso heurístico |
| `-rho` | `0.1` | Evaporación local |
| `-alpha` | `0.1` | Evaporación global |
| `-q0` | `0.9` | Prob. explotación |
| `-seed` | `0` | Semilla (0=tiempo actual) |
| `-verbose` | `false` | Salida detallada |

## Ejemplo

```bash
# Ejecutar con berlin52 (52 ciudades)
./heuristica -tsp ../Benchmark/berlin52.tsp -verbose

# Ejecutar con más iteraciones
./heuristica -tsp ../Benchmark/kroA100.tsp -iterations 500 -ants 20
```

## Estructura

```
├── main.go          # CLI
├── acs/             # Algoritmo ACS
│   ├── acs.go
│   ├── ant.go
│   ├── pheromone.go
│   └── params.go
└── tsp/             # Estructuras TSP
    ├── instance.go
    ├── tsplib.go    # Loader TSPLIB
    └── neighbor.go
```

## Benchmarks

El programa usa archivos TSPLIB ubicados en `../Benchmark/`. Incluye óptimos conocidos para comparación automática.
