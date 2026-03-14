#!/bin/bash
# ============================================================
# Script para correr todos los benchmarks TSPLIB con el AG
# Genera un archivo CSV con resultados y datos de convergencia
# ============================================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BENCHMARK_DIR="${SCRIPT_DIR}/../Benchmark"
BINARY="${SCRIPT_DIR}/tsp-ga"
OUTPUT="${SCRIPT_DIR}/resultados_ga.csv"

# Parametros del AG (modificar aqui o pasar como variables de entorno)
POP="${POP:-600}"
GEN="${GEN:-2000}"
MUT="${MUT:-0.3}"
TOURN="${TOURN:-3}"
STAG="${STAG:-200}"

# Compilar si no existe el binario o si el codigo es mas nuevo
if [ ! -f "$BINARY" ] || [ "$SCRIPT_DIR/main.go" -nt "$BINARY" ]; then
    echo "Compilando..."
    cd "$SCRIPT_DIR" && go build -o tsp-ga . || { echo "ERROR: Fallo la compilacion"; exit 1; }
fi

# Encabezado CSV
echo "benchmark,costo,tiempo,optimo,gap,pop,gen,mut,tourn,stag,gen_ultima_mejora,gen_parada,razon_parada" > "$OUTPUT"

# Contar archivos
TOTAL=$(ls "$BENCHMARK_DIR"/*.tsp 2>/dev/null | wc -l)
CURRENT=0

echo "Ejecutando $TOTAL benchmarks con: Pop=$POP Gen=$GEN Mut=$MUT Tourn=$TOURN Stag=$STAG"
echo "Resultados en: $OUTPUT"
echo ""

for TSP_FILE in "$BENCHMARK_DIR"/*.tsp; do
    CURRENT=$((CURRENT + 1))
    NOMBRE=$(basename "$TSP_FILE")
    printf "[%2d/%2d] %-20s " "$CURRENT" "$TOTAL" "$NOMBRE"

    RESULT=$("$BINARY" -flat -pop "$POP" -gen "$GEN" -mut "$MUT" -tourn "$TOURN" -stag "$STAG" "$TSP_FILE" 2>&1)

    if [ $? -eq 0 ]; then
        echo "$RESULT" >> "$OUTPUT"
        # Extraer GAP y razon para mostrar en terminal
        GAP=$(echo "$RESULT" | cut -d',' -f5)
        RAZON=$(echo "$RESULT" | cut -d',' -f13)
        GEN_MEJORA=$(echo "$RESULT" | cut -d',' -f11)
        GEN_TOTAL=$(echo "$RESULT" | cut -d',' -f12)
        printf "GAP=%6s%%  gen_mejora=%-4s  parada=%-4s  (%s)\n" "$GAP" "$GEN_MEJORA" "$GEN_TOTAL" "$RAZON"
    else
        echo "ERROR"
        echo "# ERROR: $NOMBRE - $RESULT" >> "$OUTPUT"
    fi
done

echo ""
echo "Listo. Resultados guardados en: $OUTPUT"
