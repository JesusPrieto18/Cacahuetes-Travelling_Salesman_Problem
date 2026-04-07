# Optimización por Florecimiento de Plancton (OFP) para TSP

Implementación de la metaheurística bioinspirada **Optimización por Florecimiento de Plancton (OFP)** para resolver el Problema del Viajante (Travelling Salesman Problem) usando instancias TSPLIB.

El algoritmo simula el ciclo de vida, desplazamiento y reproducción del fitoplancton en el océano, balanceando la exploración global y la explotación local.

## Requisitos

- Go 1.21 o superior
- Archivos `.tsp` en formato TSPLIB (en la carpeta `../Benchmark/`)

## Compilación y Ejecución

```bash
# Compilar el binario
go build -o tsp-ofp .

# Ejecutar con el benchmark por defecto (berlin52.tsp)
./tsp-ofp

# Ejecutar con un benchmark específico
./tsp-ofp ../Benchmark/ch150.tsp
```

### Ejecuciones Masivas

```bash
# Dar permisos de ejecución al script
chmod +x run_benchmarks.sh

# Ejecutar todas las instancias con parámetros por defecto
./run_benchmarks.sh

# Ejecutar con parámetros personalizados vía variables de entorno
POP=100 ITER=1000 DELTA=1500 GAMMA=0.99 BLOOM=0.20 ./run_benchmarks.sh
```

El script genera el archivo `resultados_ofp.csv`.

## Parámetros de la OFP

| Flag      | Tipo    | Default | Descripción                                                |
|-----------|---------|---------|------------------------------------------------------------|
| `-pop`    | int     | 50      | Tamaño de la población base ($N$)                          |
| `-iter`   | int     | 1000    | Número máximo de iteraciones                               |
| `-alpha`  | float64 | 0.1     | Intensidad de corrientes para la Deriva ($\alpha$)         |
| `-delta`  | float64 | 50.0    | Paso quimiotáctico inicial / evaluaciones 2-opt ($\delta$) |
| `-gamma`  | float64 | 0.95    | Factor de enfriamiento quimiotáctico ($\Gamma$)            |
| `-bloom`  | float64 | 0.1