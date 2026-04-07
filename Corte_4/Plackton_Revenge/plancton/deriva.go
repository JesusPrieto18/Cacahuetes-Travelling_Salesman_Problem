package plancton

import (
	"math/rand"
)

// AplicarDeriva recibe un puntero a un Plancton y modifica su ruta simulando corrientes.
// Utiliza la estrategia de desplazamiento de bloque continuo.
func AplicarDeriva(p *Plancton, alpha float64) {
	n := len(p.Tour)
	// Si el tour es muy pequeño o no hay corriente, no hacemos nada
	if n < 4 || alpha <= 0.0 {
		return
	}

	// 1. Determinar el tamaño del bloque k = ⌊α · n⌋
	k := int(alpha * float64(n))
	
	// Control de límites: El bloque debe tener entre 2 y (n-2) ciudades
	if k < 2 {
		k = 2 
	}
	if k >= n-1 {
		k = n - 2
	}

	// 2. Extraer el bloque arrastrado (manejando el arreglo como circular)
	inicio := rand.Intn(n)
	bloque := make([]int, 0, k)
	resto := make([]int, 0, n-k)

	for i := 0; i < n; i++ {
		idx := (inicio + i) % n
		if i < k {
			bloque = append(bloque, p.Tour[idx])
		} else {
			resto = append(resto, p.Tour[idx])
		}
	}

	// 3. Elegir un nuevo punto de inserción en el arreglo 'resto'
	puntoInsercion := rand.Intn(len(resto) + 1)
	
	// 4. Ensamblar el nuevo tour arrastrado por la corriente
	nuevoTour := make([]int, 0, n)
	nuevoTour = append(nuevoTour, resto[:puntoInsercion]...)
	nuevoTour = append(nuevoTour, bloque...)
	nuevoTour = append(nuevoTour, resto[puntoInsercion:]...)

	// 5. Actualizamos directamente la estructura original
	p.Tour = nuevoTour
	
	// Nota de diseño: No actualizamos p.Cost aquí intencionalmente.
	// Físicamente el plancton es arrastrado sin evaluar el entorno.
	// La siguiente fase (Quimiotaxis) se encargará de evaluar y actualizar el fitness.
}