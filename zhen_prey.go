package main

// Input: currentUnit is a pointer to a unit, currentEcosystem is a pointer to the ecosystem, i and j are the indices of the location of the unit we are about to move, curGen is the number of generations of the unit we are about to move.
// Output: none, operates on pointers
func MovePrey(currentUnit *Unit, currentEcosystem *Ecosystem, i, j, curGen int) {
	currentPrey := currentUnit.prey
	deltaX, deltaY, moveIndex := UseGenomeToMove(currentPrey)
	//movementDeltas store orderedPair of deltaX and deltaY
	movementDeltas = deltas[currentPrey.lastDirection]
	currentUnit.prey = nil

	// check if the prey eats
	//if (*currentEcosystem)[i+deltaX][j+deltaY].food.isPresent == true {
	if CheckIfEats(currentEcosystem[i+deltaX][j+deltaY], currentPrey) {
		currentPrey.FeedOrganism()
	}

	// energy decreases based on how drastic the change in direction is for the movement
	currentPrey.DecreaseEnergy()

	// check if there is something in the unit already. if there isn't anything then the prey doesn't move. alternatively we can randomly choose an open Unit adjacnet to the prey
	if (*currentEcosystem)[i+deltaX][j+deltaY].prey == nil && (*currentEcosystem)[i+deltaX][j+deltaY].predator == nil {
		(*currentEcosystem)[i+deltaX][j+deltaY].prey = currentPrey
		currentPrey.lastDirection = lastDirection
		currentPrey.lastGenUpdated = curGen
	}

}

func CheckIfEats(currUnit *Unit, currPrey *Prey) bool {
	return currUnit.food.isPresent && currPrey.energy < maxEnergy
}

func (currPrey Prey) FeedOrganism(currUnit *Unit) {
	currUnit.food.isPresent = false
	currPrey.energy += 1
}
//Set a constant dictionary where keys are the directionIndex and the values are the orderedPair with corresponding deltaX and deltaY
deltas := map[int]OrderedPair {
	0: OrderedPair{-1, 1},
	1: OrderedPair{0, 1},
	2: OrderedPair{1, 1},
	3: OrderedPair{-1, 0},
	4: OrderedPair{1, 0},
	5: OrderedPair{-1, -1},
	6: OrderedPair{0, -1},
	7: OrderedPair{1, -1},
}

// cannot move to unit where there's shark (predator)
//deltaX, deltaY, lastDirection := UseGenomeToMove(currentPrey)
func UseGenomeToMove(*currPrey) int, int {
	r := rand.Float64()
	directionIndex := 0
	runningSum := 0
	for idx, gene := range currPrey.genome {
		runningSum += gene
		if runningSum >= r {
			directionIndex = idx
			break
		}
	}
	//lastDirection will be updated with my new direction
	currentPrey.lastDirection = (currentPrey.lastDirection + directionIndex) % 8
	moveDeltas := deltas[currentPrey.lastDirection]
	return moveDeltas.x, moveDeltas.y
}

func (currPrey Prey) DecreaseEnergy(currPrey *Prey) {
	currPrey.energy -= 1
	if currPrey.lastDirection == 1 || currPrey.lastDirection == 7 {
		currPrey.energy -= 1
	} else if currPrey.lastDirection == 2 || currPrey.lastDirection == 6 {
		currPrey.energy -= 2
	} else if currPrey.lastDirection == 3 || currPrey.lastDirection == 5 {
		currPrey.energy -= 4
	} else if currPrey.lastDirection == 4 {
		currPrey.energy -= 8
	}
}
