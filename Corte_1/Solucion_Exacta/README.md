# TSP Solver: Branch and Bound con Lower Bound

Se implementa una solución exacta para el Problema del Viajero (TSP) utilizando el algoritmo Branch and Bound (Ramificación y Poda).


Este algoritmo resuelve el TSP de forma exacta mediante Best-First Search y una cola de prioridad (Min-Heap) que prioriza nodos con menor cota inferior (lowerBound). El algoritmo utiliza poda (pruning), que descarta ramas cuyo lowerBound supera al mejor costo hallado, basando este cálculo en la suma de las dos aristas más baratas por ciudad y el ajuste del camino fijo.

## Requisitos y Compilación

1.  Requiere Go versión 1.21 o superior.
2.  Para compilar el proyecto, ejecute:
    ```bash
    go build -o tsp_solver main.go
    ```

## Ejecución

Para ejecutar se tiene que pasar la ruta del archivo de la instancia como argumento:

```bash
./tsp_solver -tsp ruta/al/archivo.tsp