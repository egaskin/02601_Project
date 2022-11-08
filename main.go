package main

func main() {

	var initialEcosystem Ecosystem = InitializeEcosystem()
	var numGens int = 100000
	var foodRule string = "eden"
	var predator1 Predator

	predator1.Organism.age = 10
	predator1.age = 10

	allEcosystems := SimulateEcosystemEvolution(&initialEcosystem, numGens, foodRule)

}
