#!/bin/bash
# ============================================================
# Script para correr todos los benchmarks TSPLIB con la OFP
# Genera un archivo CSV con resultados y datos de convergencia
# ============================================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BENCHMARK_DIR="${SCRIPT_DIR}/../Benchmark"
BINARY="${SCRIPT_DIR}/tsp-ofp"
OUTPUT="${SCRIPT_DIR}/resultados_ofp.csv"

# Parametros de la OFP (modificar aqui o pasar como variables de entorno)
POP="${POP:-100}"
ITER="${ITER:-1000}"
ALPHA="${ALPHA:-0.10}"
DELTA="${DELTA:-1500.0}"
GAMMA="${GAMMA:-0.99}"
BLOOM="${BLOOM:-0.20}"
TFREQ="${TFREQ:-50}"
TMU="${TMU:-0.20}"

# Compilar si no existe el binario o si el codigo es mas nuevo
if [ ! -f "$BINARY" ] || [ "$SCRIPT_DIR/main.go" -nt "$BINARY" ]; then
    echo "Compilando OFP..."
    cd "$SCRIPT_DIR" && go build -o tsp-ofp . || { echo "ERROR: Falló la compilación"; exit 1; }
fi

# Encabezado CSV adaptado a la salida plana de la OFP
echo "benchmark,costo,tiempo,optimo,gap,pop,iter,alpha,delta,gamma,bloom,tfreq,tmu,gen_ultima_mejora" > "$OUTPUT"

# Contar archivos
TOTAL=$(ls "$BENCHMARK_DIR"/*.tsp 2>/dev/null | wc -l)
CURRENT=0

echo "Ejecutando $TOTAL benchmarks con OFP:"
echo "Pop=$POP Iter=$ITER Alpha=$ALPHA Delta=$DELTA Gamma=$GAMMA Bloom=$BLOOM TFreq=$TFREQ TMu=$TMU"
echo "Resultados en: $OUTPUT"
echo ""

for TSP_FILE in "$BENCHMARK_DIR"/*.tsp; do
    CURRENT=$((CURRENT + 1))
    NOMBRE=$(basename "$TSP_FILE")
    printf "[%2d/%2d] %-20s " "$CURRENT" "$TOTAL" "$NOMBRE"

    # Ejecutar binario
    RESULT=$("$BINARY" -flat -pop "$POP" -iter "$ITER" -alpha "$ALPHA" -delta "$DELTA" -gamma "$GAMMA" -bloom "$BLOOM" -tfreq "$TFREQ" -tmu "$TMU" "$TSP_FILE" 2>&1)

    if [ $? -eq 0 ]; then
        echo "$RESULT" >> "$OUTPUT"
        
        # Extraer variables para mostrar en consola (basado en las columnas del flat format)
        # Archivo(1),Costo(2),Tiempo(3),Optimo(4),GAP(5),Pop(6),Iter(7),Alpha(8),Delta(9),Gamma(10),Bloom(11),TFreq(12),TMu(13),UltimaMejora(14)
        GAP=$(echo "$RESULT" | cut -d',' -f5)
        TIEMPO=$(echo "$RESULT" | cut -d',' -f3)
        GEN_MEJORA=$(echo "$RESULT" | cut -d',' -f14)
        
        printf "GAP=%6s%%  Tiempo=%-12s  gen_mejora=%s\n" "$GAP" "$TIEMPO" "$GEN_MEJORA"
    else
        echo "ERROR"
        echo "# ERROR: $NOMBRE - $RESULT" >> "$OUTPUT"
    fi
done

echo ""
echo "Listo. Resultados guardados en: $OUTPUT"