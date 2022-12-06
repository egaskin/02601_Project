package main

import "math/rand"

//Set a constant dictionary where keys are the directionIndex and the values are the orderedPair with corresponding deltaX and deltaY
// movementDeltas := map[int]OrderedPair {
// 	0: OrderedPair{-1, -1},
// 	1: OrderedPair{0, -1},
// 	2: OrderedPair{1, -1},
// 	3: OrderedPair{1, 0},
// 	4: OrderedPair{1, 1},
// 	5: OrderedPair{0, 1},
// 	6: OrderedPair{-1, 1},
// 	7: OrderedPair{-1, 0},
// }

// Gene index to energy cost
// energyCosts := map[int]int {
// 	0: 0,
// 	1: -1,
// 	2: -2,
// 	3: -4,
// 	4: -8,
// 	5: -4,
// 	6: -2,
// 	7: -1,
// }

// Input: currentUnit is a pointer to a unit, currentEcosystem is a pointer to the ecosystem, i and j are the indices of the location of the unit we are about to move, curGen is the number of generations of the unit we are about to move.
// Output: none, operates on pointers
func MovePrey(currentUnit *Unit, currentEcosystem *Ecosystem, i, j, curGen int) {
	currentPrey := currentUnit.prey
	deltaX, deltaY, newDirection, geneIndex, newI, newJ := UseGenomeToMove(currentEcosystem, currentPrey, i, j)

	// check if the prey eats
	//if (*currentEcosystem)[i+deltaX][j+deltaY].food.isPresent == true {
	// if CheckIfEats((*currentEcosystem)[i+deltaX][j+deltaY], currentPrey) {
	// 	currentPrey.FeedOrganism((*currentEcosystem)[i+deltaX][j+deltaY])
	// }

	if CheckIfEats((*currentEcosystem)[newI][newJ], currentPrey) {
		currentPrey.FeedOrganism((*currentEcosystem)[newI][newJ])
	}

	// energy decreases based on how drastic the change in direction is for the movement
	// if at least one of deltaX or deltaY is not equal to 0, we move the prey
	isMoving := deltaX != 0 || deltaY != 0
	currentPrey.DecreaseEnergy(geneIndex, isMoving)

	currentUnit.prey = nil

	// check if energy level > 0
	// if it is not, update direction
	if currentPrey.energy > 0 {
		currentPrey.lastDirection = newDirection
		// when deltaX and deltaY == 0, currentPrey stay at unit [i, j]
		(*currentEcosystem)[i+deltaX][j+deltaY].prey = currentPrey
	}

}

func CheckIfEats(currentUnit *Unit, currentPrey *Prey) bool {
	return currentUnit.food.isPresent && (currentPrey.energy < maxEnergy)
}

func (currentPrey Prey) FeedOrganism(currentUnit *Unit) {
	currentUnit.food.isPresent = false
	currentPrey.energy += 1
}

// cannot move to unit where there's shark (predator)
// deltaX, deltaY, lastDirection := UseGenomeToMove(currentPrey)
func UseGenomeToMove(currentEcosystem *Ecosystem, currentPrey *Prey, i, j int) (int, int, int, int, int, int) {
	var moveDeltas OrderedPair
	var geneIndex, newDirection, newI, newJ int
	isFreeUnitFlag := false
	numTries := 0

	// 20 is the threshold for max number of tries we get to reselect a gene for movement
	// if numberTries >= 20 and isFreeUnitFlag is still false
	// the prey doesn't move
	for !isFreeUnitFlag && numTries < 20 {
		r := rand.Float64()
		geneIndex = 0
		runningSum := 0.0
		for idx, gene := range currentPrey.genome {
			runningSum += float64(gene)
			if runningSum >= r {
				geneIndex = idx
				break
			}
		}
		newDirection := (currentPrey.lastDirection + geneIndex) % 8
		moveDeltas = deltas[newDirection]
		numRows := len(*currentEcosystem)
		numCols := len((*currentEcosystem)[0])
		newI = GetIndex(i, moveDeltas.row, numRows)
		newJ = GetIndex(j, moveDeltas.col, numCols)

		isFreeUnitFlag = isFreeUnit(currentEcosystem, i, j)
		numTries += 1
	}
	// if numTries >= 20 and still haven't find a free unit, we don't move
	if !isFreeUnitFlag {
		geneIndex = 0
		newDirection = currentPrey.lastDirection
		moveDeltas.row, moveDeltas.col = 0, 0
	}

	//lastDirection will be updated with my new direction
	return moveDeltas.row, moveDeltas.col, newDirection, geneIndex, newI, newJ
}

func (currentPrey *Prey) DecreaseEnergy(geneIndex int, isMoving bool) {
	currentPrey.energy -= 1

	// if prey needs to be moved since either deltaX or deltaY or both are not equal to 0
	// we decrease the energy based on the geneIndex
	if isMoving {
		currentPrey.energy -= energyCosts[geneIndex]
	}
}

// check if unit (i, j) is unoccupied by a predator or another prey
func isFreeUnit(currentEcosystem *Ecosystem, i, j int) bool {
	if (*currentEcosystem)[i][j].prey == nil && (*currentEcosystem)[i][j].predator == nil {
		return true
	} else {
		return false
	}
}

// pass in the row, column indices and the delta for movement
// return new row, column indices within the boundary
// boundary is the numRow and numCol of the ecosystem board
func GetIndex(index, delta, boundary int) int {
	newIndex := index + delta
	if newIndex < 0 {
		newIndex = boundary + newIndex
	} else if newIndex >= boundary {
		newIndex = newIndex % boundary
	}
	return newIndex
}
