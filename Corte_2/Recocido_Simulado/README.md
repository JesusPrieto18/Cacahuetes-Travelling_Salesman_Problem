# Recocido Simulado - TSP

Este programa implementa el algoritmo de Recocido Simulado (Simulated Annealing) para resolver el Problema del Viajante (TSP), permitiendo ajustar parámetros desde la línea de comandos.

Este programa implementa el algoritmo de Recocido Simulado (Simulated Annealing) para resolver el Problema del Viajante (TSP).

## Requisitos
- Go 1.18 o superior

## Estructura del Proyecto
- `main.go`: Archivo principal para ejecutar el programa.
- `localsearch/`: Búsqueda local (2-opt).
- `models/`: Definición de estructuras de datos.
- `parser/`: Lectura de archivos TSP.
- `solver/`: Lógica de solución (incluye Recocido Simulado).
- `utils/`: Utilidades varias.
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

3. Ejecuta el programa indicando la ruta del archivo TSP y parámetros opcionales:

```
recocido_simulado.exe ../Benchmark/berlin52.tsp -temp=1000 -alpha=0.995 -min_temp=0.001 -iter=1000
```

O directamente con `go run`:

```
go run main.go ../Benchmark/berlin52.tsp -temp=1000 -alpha=0.995 -min_temp=0.001 -iter=1000
```

## Parámetros por línea de comandos
- `-temp`: Temperatura inicial del recocido (float).
- `-alpha`: Factor de enfriamiento (float, <1).
- `-min_temp`: Temperatura mínima de parada (float).
- `-iter`: Iteraciones por nivel de temperatura (int).

## Ejemplo de salida
El programa mostrará en consola la mejor ruta encontrada, su costo total, el óptimo (si está disponible) y el GAP:

```
Configuración SA: Temp=1000.00, Alpha=0.9950, Min=0.0010, Iter=1000
Benchmark: berlin52.tsp
Tiempo     	Costo     	Optimo	GAP SA (%)
0.123456s	7542.0	    7542	0.00
```
