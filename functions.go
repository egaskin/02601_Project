package main

import "fmt"

func SimulateEcosystemEvolution(initialEcosystem *Ecosystem, numGens int, foodRule string) []*Ecosystem {
	var allEcosystems []*Ecosystem
	allEcosystems[0] = initialEcosystem

	for i := 1; i < numGens; i++ {
		allEcosystems[i] = UpdateEcosystem(allEcosystems[i-1], foodRule, i)
	}

	return allEcosystems
}

func UpdateEcosystem(prevEcosystem *Ecosystem, foodRule string, curGen int) *Ecosystem {
	// var row int = len(*prevEcosystem)
	// var col int = len((*prevEcosystem)[0])
	var nextEcosystem Ecosystem = DeepCopyEcosystem(prevEcosystem)

	for i := range nextEcosystem {
		for j := range nextEcosystem[i] {
			// get a pointer to the current Unit that we need new values for (nextEcosystem)
			currentUnit := nextEcosystem[i][j]

			fmt.Println("SOMETHING ABOUT RANDOMLY CHOOSING AN UNCHOSEN i here")

			// Update the currentUnit based on the nextEcosystem! since we want the system to change as prey and food or disappearing (so each prey/predator is competing to get to their respective food source first)

			// only perform this operation is Unit contains predator
			if (*currentUnit).predator != nil {

				// skip already updated predator.
				if (*currentUnit).predator.lastGenUpdated != curGen {
					currentUnit.predator.UpdatePred(nextEcosystem, i, j)
					currentUnit.predator.lastGenUpdated = curGen
				}

			} else if (*currentUnit).prey != nil { // only perform this operation if Unit contains predator

				// skip already updated prey.
				if (*currentUnit).prey.lastGenUpdated != curGen {
					currentUnit.prey.UpdatePrey(nextEcosystem, i, j)
					currentUnit.prey.lastGenUpdated = curGen
				}

			} else if (*currentUnit).food.isPresent != true { // skip if the food is already true. note: we don't need to check the lastGenUpdated because food will be false if this is ran.

				// determine whether food appears randomly for the prey. GeneratePreyFoodRandomly() will update both fields of the food, if food is generated. otherwise it will leave it false.
				currentUnit.food.GeneratePreyFoodRandomly(foodRule, curGen, i, j)
			}
		}
	}

	return &nextEcosystem
}
