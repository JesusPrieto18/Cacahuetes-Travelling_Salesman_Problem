# Algoritmo Genetico para TSP

Implementacion de un Algoritmo Genetico (AG) para resolver el Problema del Viajante (Travelling Salesman Problem) usando instancias TSPLIB.

## Requisitos

- Go 1.21 o superior
- Archivos `.tsp` en formato TSPLIB (carpeta `../Benchmark/`)

## Compilacion y Ejecucion

```bash
# Compilar
cd Corte_3/Busqueda_Dispersa
go build -o tsp-ds .

# Ejecutar con benchmark por defecto (berlin52)
./tsp-ds

# Ejecutar con un benchmark especifico
./tsp-ds ../Benchmark/eil76.tsp

# Ejecutar sin compilar
go run . ../Benchmark/kroA100.tsp
```

### Benchmarks masivos

```bash
# Ejecutar todas las instancias con parametros por defecto
./run_benchmarks.sh

# Ejecutar con parametros personalizados via variables de entorno
POP=1000 GEN=800 MUT=0.25 TOURN=5 STAG=150 RELINK=0.4 DIVTHRESH=7 ./run_benchmarks.sh
```

El script genera el archivo `resultados_ds.csv`.

## Parametros

| Flag    | Tipo    | Default | Descripcion                                              |
|---------|---------|---------|----------------------------------------------------------|
| `-pop`  | int     | 600     | Tamaño de la poblacion                                   |
| `-gen`  | int     | 2000    | Numero maximo de generaciones                            |
| `-mut`  | float64 | 0.3     | Probabilidad de mutacion (0.0 a 1.0)                    |
| `-tourn`| int     | 3       | Tamaño del torneo para seleccion de padres               |
| `-stag` | int     | 200     | Generaciones sin mejora antes de parar (0 = desactivado) |
| `-relink` | float64 | 0.5   | Porcentaje de pares de soluciones a reenlazar por generacion |
| `-divthresh` | int | 5      | Distancia minima en aristas para aceptar un individuo     |
| `-flat` | bool    | false   | Salida en formato plano separado por comas (sin encabezados) |

### Ejemplos

```bash
# Configuracion personalizada
./tsp-ds -pop 1000 -gen 800 -mut 0.25 -tourn 5 -relink 0.4 -divthresh 7 ../Benchmark/berlin52.tsp

# Desactivar parada por estancamiento
./tsp-ds -stag 0 ../Benchmark/eil76.tsp

# Ajustar intensidad del relinking y filtro de diversidad
./tsp-ds -relink 0.6 -divthresh 10 ../Benchmark/kroA100.tsp

# Salida plana (util para scripts y pipelines)
./tsp-ds -flat ../Benchmark/kroA100.tsp
```

## Componentes del Algoritmo

### 1. Representacion

- **Genotipo:** Permutacion de indices `[0..N-1]` donde cada indice representa una ciudad.
- **Fenotipo:** Tour (circuito hamiltoniano) obtenido al recorrer las ciudades en el orden de la permutacion.

### 2. Poblacion Inicial

La poblacion se inicializa con tres tipos de individuos para garantizar diversidad y calidad:

1. **Farthest Insertion (1 individuo):** Heuristica constructiva que genera un tour de buena calidad. Parte de las dos ciudades mas lejanas y en cada paso inserta la ciudad mas lejana al tour en la posicion de menor costo.
2. **Variantes perturbadas (~15% de la poblacion):** Copias del tour de Farthest Insertion con swaps aleatorios aplicados, para tener individuos buenos pero diferentes entre si.
3. **Permutaciones aleatorias (~85%):** El resto se genera aleatoriamente para mantener diversidad genetica.

Se aplica **control de diversidad**: individuos con costos duplicados se descartan y se regeneran.

### 3. Seleccion de Padres — Torneo

Se seleccionan `k` individuos al azar de la poblacion y el de menor costo (mejor aptitud) se elige como padre. Los individuos peores tienen probabilidad positiva de ser seleccionados si compiten contra otros aun peores, lo que mantiene diversidad.

### 4. Cruce — Corte y Llenado (Order Crossover)

Operador descrito en Clase 6:

1. Se elige un punto de corte `p` aleatorio.
2. **Hijo 1:** Copia el prefijo `padre1[0..p)`, luego recorre `padre2` en orden y agrega los elementos que no estan ya en el hijo.
3. **Hijo 2:** Simetrico — prefijo de `padre2`, completar con `padre1`.

Esto garantiza que los hijos sean permutaciones validas.

### 5. Mutacion — Inversion

Se escogen dos posiciones aleatorias `i` y `j` y se invierte el segmento entre ellas. Este operador es efectivo para TSP porque preserva la adyacencia de la mayoria de las ciudades.

Se aplica con probabilidad configurable (`-mut`).

### 6. Reenlace de Caminos (Path Relinking)

La busqueda dispersa incorpora una fase de **path relinking** configurable con `-relink`.
En cada generacion se toma un porcentaje de pares de soluciones y se exploran trayectorias entre ellas para construir individuos intermedios de alta calidad.

Esta fase intensifica la busqueda en regiones prometedoras del espacio de soluciones sin abandonar por completo la diversidad poblacional.

### 7. Control de Diversidad

Para evitar convergencia prematura, un individuo nuevo solo se acepta si mantiene una distancia minima respecto a la poblacion actual.
Ese umbral se controla con `-divthresh` y se mide en diferencias de aristas del tour.

Valores bajos favorecen explotacion; valores mas altos favorecen exploracion.

### 8. Seleccion de Sobrevivientes — (μ+λ)

En cada generacion se unen la poblacion actual (μ) con los hijos generados (λ), se ordenan por costo y se seleccionan los mejores `PopSize` individuos. Esto es mas robusto que el reemplazo generacional puro porque nunca pierde buenas soluciones.

### 9. Criterios de Terminacion

El algoritmo se detiene cuando se cumple **cualquiera** de estas condiciones:

1. **Maximo de generaciones:** Se alcanza el limite configurado con `-gen`.
2. **Estancamiento:** El mejor costo no mejora durante `-stag` generaciones consecutivas.

## Estructura del Proyecto

```
Busqueda_Dispersa/
├── main.go                        # Punto de entrada, flags y salida
├── go.mod                         # Modulo Go
├── models/
│   └── city.go                    # Struct City (ID, X, Y)
├── parser/
│   └── reader.go                  # Lector de archivos .tsp (formato TSPLIB)
├── geneticalgorithm/
│   ├── ga.go                      # Logica principal del AG (poblacion, bucle, reemplazo)
│   ├── crossover.go               # Cruce Corte y Llenado (Order Crossover)
│   ├── mutation.go                # Mutacion por Inversion
│   ├── relinking.go               # Path relinking entre soluciones prometedoras
│   ├── selection.go               # Seleccion por Torneo
│   └── heuristic.go               # Heuristica Farthest Insertion
├── solver/
│   └── ga_solver.go               # Wrapper del solver
└── utils/
    ├── utils.go                   # Distancia euclidiana, costo total
    └── bks.go                     # Best Known Solutions de TSPLIB
```

## Salida

### Formato normal (por defecto)

```
Benchmark   Tiempo      Costo       Optimo  GAP AM (%)
berlin52.tsp  143ms     7701.4556   7542    2.11
Configuracion AM: Pop=600, Gen=2000, Mut=0.3000, Tourn=3, Stag=200, Relink=0.50, DivThresh=5
Convergencia: ultima mejora en gen 184, parada en gen 384 por stagnation_limit
```

### Formato plano (`-flat`)

Columnas separadas por comas, sin encabezados (una linea):

```
Archivo,Costo,Tiempo,Optimo,GAP,Pop,Gen,Mut,Tourn,Stag,Relink,DivThresh,GenUltimaMejora,GenParada,RazonParada
```
