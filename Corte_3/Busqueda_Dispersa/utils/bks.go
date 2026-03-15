package utils

import (
	"path/filepath"
	"strings"
)

// TSPLIBOptimal contains known optimal solutions for available benchmarks
var TSPLIBOptimal = map[string]float64{
	"berlin52":  7542,
	"bier127":   118282,
	"brd14051":  469385,
	"ch130":     6110,
	"ch150":     6528,
	"d198":      15780,
	"d493":      35002,
	"d657":      48912,
	"d1291":     50801,
	"d1655":     62128,
	"d2103":     80450,
	"d15112":    1573084,
	"eil51":     426,
	"eil76":     538,
	"eil101":    629,
	"fl417":     11861,
	"fl1400":    20127,
	"fl1577":    22249,
	"fl3795":    28772,
	"fnl4461":   182566,
	"gil262":    2378,
	"kroA100":   21282,
	"kroA150":   26524,
	"kroA200":   29368,
	"kroB100":   22141,
	"kroB150":   26130,
	"kroC100":   20749,
	"kroD100":   21294,
	"kroE100":   22068,
	"lin105":    14379,
	"lin318":    42029,
	"nrw1379":   56638,
	"p654":      34643,
	"pcb442":    50778,
	"pcb1173":   56892,
	"pcb3038":   137694,
	"pr76":      108159,
	"pr107":     44303,
	"pr124":     59030,
	"pr136":     96772,
	"pr144":     58537,
	"pr152":     73682,
	"pr226":     80369,
	"pr264":     49135,
	"pr299":     48191,
	"pr439":     107217,
	"pr1002":    259045,
	"pr2392":    378032,
	"rat99":     1211,
	"rat195":    2323,
	"rat575":    6773,
	"rat783":    8806,
	"rd100":     7910,
	"rd400":     15281,
	"rl1304":    252948,
	"rl1323":    270199,
	"rl1889":    316536,
	"rl5915":    565530,
	"rl5934":    556045,
	"rl11849":   923288,
	"st70":      675,
	"ts225":     126643,
	"tsp225":    3916,
	"u159":      42080,
	"u574":      36905,
	"u724":      41910,
	"u1060":     224094,
	"u1432":     152970,
	"u1817":     57201,
	"u2152":     64253,
	"u2319":     234256,
	"vm1084":    239297,
	"vm1748":    336556,
}

// GetOptimalCost intenta obtener el costo óptimo basado en el nombre del archivo
func GetOptimalCost(filename string) float64 {
	// 1. Obtener el nombre base (ej: "../Benchmark/berlin52.tsp" -> "berlin52.tsp")
	base := filepath.Base(filename)

	// 2. Quitar la extensión (ej: "berlin52.tsp" -> "berlin52")
	name := strings.TrimSuffix(base, filepath.Ext(base))

	// 3. Buscar en el mapa
	if val, ok := TSPLIBOptimal[name]; ok {
		return val
	}
	return 0 // Retorna 0 si no se encuentra
}
