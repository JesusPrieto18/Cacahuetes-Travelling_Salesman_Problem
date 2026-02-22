package geneticalgorithm

import "math/rand"

// CutAndFillCrossover implements the "corte y llenado" (Order Crossover) operator.
// Selects a random cut point p. Child1 takes parent1[0..p], then fills with
// elements from parent2 in order, skipping those already present.
// Child2 is symmetric: prefix from parent2, fill from parent1.
func CutAndFillCrossover(parent1, parent2 []int) ([]int, []int) {
	n := len(parent1)

	// Cut point: at least 1 element from prefix, at least 1 to fill
	p := rand.Intn(n-1) + 1 // p in [1, n-1)

	child1 := cutAndFillBuild(parent1, parent2, p, n)
	child2 := cutAndFillBuild(parent2, parent1, p, n)

	return child1, child2
}

// cutAndFillBuild builds one child: copies prefix[0..p) from donor,
// then appends elements from filler in order, skipping duplicates.
func cutAndFillBuild(donor, filler []int, p, n int) []int {
	child := make([]int, 0, n)
	inChild := make([]bool, n)

	// Copy prefix from donor
	for i := 0; i < p; i++ {
		child = append(child, donor[i])
		inChild[donor[i]] = true
	}

	// Fill remaining from filler in order
	for _, val := range filler {
		if !inChild[val] {
			child = append(child, val)
			inChild[val] = true
		}
	}

	return child
}
