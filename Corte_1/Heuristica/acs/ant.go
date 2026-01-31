package acs

// Ant represents a single ant in the colony
type Ant struct {
	tour      []int  // current tour being constructed
	visited   []bool // cities already visited
	tourIndex int    // current position in tour construction
	n         int    // number of cities
}

// NewAnt creates a new ant for a problem with n cities
func NewAnt(n int) *Ant {
	return &Ant{
		tour:    make([]int, n),
		visited: make([]bool, n),
		n:       n,
	}
}

// Reset prepares the ant for a new tour starting from the given city
func (a *Ant) Reset(startCity int) {
	for i := range a.visited {
		a.visited[i] = false
	}
	a.tour[0] = startCity
	a.visited[startCity] = true
	a.tourIndex = 1
}

// CurrentCity returns the city where the ant is currently located
func (a *Ant) CurrentCity() int {
	if a.tourIndex == 0 {
		return a.tour[0]
	}
	return a.tour[a.tourIndex-1]
}

// CanVisit returns true if the city has not been visited yet
func (a *Ant) CanVisit(city int) bool {
	return !a.visited[city]
}

// Visit moves the ant to the specified city
func (a *Ant) Visit(city int) {
	a.tour[a.tourIndex] = city
	a.visited[city] = true
	a.tourIndex++
}

// Tour returns the complete tour (only valid after all cities visited)
func (a *Ant) Tour() []int {
	return a.tour
}

// TourComplete returns true if the ant has visited all cities
func (a *Ant) TourComplete() bool {
	return a.tourIndex >= a.n
}

// UnvisitedCities returns a slice of cities not yet visited
func (a *Ant) UnvisitedCities() []int {
	unvisited := make([]int, 0, a.n-a.tourIndex)
	for i := 0; i < a.n; i++ {
		if !a.visited[i] {
			unvisited = append(unvisited, i)
		}
	}
	return unvisited
}
