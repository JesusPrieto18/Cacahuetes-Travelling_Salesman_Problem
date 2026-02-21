# GRASP Reactivo - TSP

Este programa implementa un algoritmo GRASP Reactivo (Greedy Randomized Adaptive Search Procedure) para resolver el Problema del Viajante (TSP).

## Requisitos
- Go 1.18 o superior

## Estructura del proyecto
- `main.go`: Archivo principal para ejecutar el programa.
- `grasp/`: Construcción GRASP Reactiva y lógica de selección de sesgo.
- `localsearch/`: Busqueda local (2-opt).
- `models/`: Definición de estructuras de datos.
- `parser/`: Lectura de archivos TSP.
- `utils/`: Utilidades varias.
- Los archivos de instancias TSP se encuentran en la carpeta `../Benchmark/`.

## Cómo compilar y ejecutar

1. Abre una terminal en la carpeta `GRASP`:

```
cd Corte_2\GRASP
```

2. Compila el programa:

```
go build -o grasp.exe main.go
```

3. Ejecuta el programa indicando la ruta del archivo TSP:

```
grasp.exe ../Benchmark/berlin52.tsp
```

O directamente con `go run`:

```
go run main.go ../Benchmark/berlin52.tsp
```

## Parámetros por linea de comandos
- La instancia TSP se pasa como primer argumento (si no se especifica, usa `../Benchmark/berlin52.tsp`).
- El numero de iteraciones del GRASP Reactivo esta definido en el código (`1000`). Si deseas modificarlo, ajusta el segundo parámetro en `grasp.GraspReactivo(cities, 1000)` dentro de `main.go`.

## Ejemplo de salida
El programa mostrara en consola el resultado con el tiempo, el costo obtenido, el óptimo y el GAP.

```
Instancia  Tiempo   Resultado  Optimo  GAP (%)
berlin52.tsp 0.123s 7542.00    7542.00 0.00
```
