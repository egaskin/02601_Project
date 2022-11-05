package main

func main() {

	var initialEcosystem Ecosystem = InitializeEcosystem()
	var numGens int = 100000
	var foodRule string = "eden"

	allEcosystem := SimulateEcosystemEvolution(&initialEcosystem, numGens, foodRule)

}
