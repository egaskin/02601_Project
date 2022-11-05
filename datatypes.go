package main

type Ecosystem [][]*Unit // is a 2D array of Unit

type Unit struct {
	food     bool
	predator *Predator
	prey     *Prey
}

type Organism struct {
	// we don't need location OrderedPair because we are using an [][]Cell
	energy int
	age    int
	genome [8]Gene
}

type Gene float64 // with range 0 to 1. all the genes of a genome add up to 1

type Prey struct {
	Organism
}

type Predator struct {
	Organism
}
