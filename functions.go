package main

func SimulateEcosystemEvolution(initialEcosystem *Ecosystem, numGens int, foodRule string) []*Ecosystem {
	var allEcosystems []*Ecosystem
	allEcosystems[0] = initialEcosystem

	for i := 1; i < numGens; i++ {
		allEcosystems[i] = UpdateEcosystem(allEcosystems[i-1], foodRule)
	}

	return allEcosystems
}

func UpdateEcosystem(prevEcosystem *Ecosystem, foodRule string) *Ecosystem {
	// var row int = len(*prevEcosystem)
	// var col int = len((*prevEcosystem)[0])
	var nextEcosystem Ecosystem = DeepCopyEcosystem(prevEcosystem)

	for i := range nextEcosystem {
		for j := range nextEcosystem[i] {
			// get a pointer to the current Unit that we need new values for (nextEcosystem)
			currentUnit := nextEcosystem[i][j]

			// Update the currentUnit based on the nextEcosystem! since we want the system to change as prey and food or disappearing (so each prey/predator is competing to get to their respective food source first)
			currentUnit.predator.UpdatePred(nextEcosystem, i, j)
			currentUnit.prey.UpdatePrey(nextEcosystem, i, j)
			currentUnit.food.GeneratePreyFood(foodRule, i, j) // determine whether food appears randomly for the pre

		}
	}

	return &nextEcosystem
}
