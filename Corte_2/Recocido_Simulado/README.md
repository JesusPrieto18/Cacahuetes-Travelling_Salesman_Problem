# Recocido Simulado - TSP

Este programa implementa el algoritmo de Recocido Simulado (Simulated Annealing) para resolver el Problema del Viajante (TSP).

## Requisitos
- Go 1.18 o superior

## Estructura del Proyecto
- `main.go`: Archivo principal para ejecutar el programa.
- `localsearch/`, `models/`, `parser/`, `solver/`, `utils/`: Carpetas con módulos y utilidades.
- Los archivos de instancias TSP se encuentran en la carpeta `../Benchmark/`.

## Cómo compilar y ejecutar

1. Abre una terminal en la carpeta `Recocido_Simulado`:

```
cd Corte_2\Recocido_Simulado
```

2. Compila el programa:

```
go build -o recocido_simulado.exe main.go
```

3. Ejecuta el programa indicando la ruta del archivo TSP:

```
recocido_simulado.exe ../Benchmark/berlin52.tsp
```

O directamente con `go run`:

```
go run main.go ../Benchmark/berlin52.tsp
```

## Parámetros
Algunos parámetros pueden estar configurados dentro del código fuente (`sa_solver.go`). Revisa y ajusta según tus necesidades:
- Temperatura inicial
- Factor de enfriamiento
- Número de iteraciones

## Salida
El programa mostrará en consola la mejor ruta encontrada y su costo total.

## Créditos
Desarrollado para la materia de Optimización de Rutas.
