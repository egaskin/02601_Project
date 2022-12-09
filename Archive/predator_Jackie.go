package main

import "math/rand"


// Set a constant dictionary where keys are the directionIndex and the values are the orderedPair with corresponding deltaX and deltaY
movementDeltas := map[int]OrderedPair {
	0: OrderedPair{-1, -1},
	1: OrderedPair{0, -1},
	2: OrderedPair{1, -1},
	3: OrderedPair{1, 0},
	4: OrderedPair{1, 1},
	5: OrderedPair{0, 1},
	6: OrderedPair{-1, 1},
	7: OrderedPair{-1, 0},
}

// Gene index to energy cost
energyCosts := map[int]int {
	0: 0,
	1: -1,
	2: -2,
	3: -4,
	4: -8,
	5: -4,
	6: -2,
	7: -1,
}


// UpdatePredator is a Predator method which will take a Predator input and update the position, initiate eating, reproduction, and age accordingly
func (shark *Predator) UpdatePredator(currEco *Ecosystem, i, j int) {

	if shark.Organism.energy == 0 {

		(*currEco)[i][j].predator = nil

	} else {

		//1. Update POSITION AND ENERGY first if energy is allowed
		// This UpdatePosition will scan through all 7 units, give a list of available units, and use GENOME to update
		//	We prioritize the GENOME instead of the fish
		// This function will UpdatePredatorPosition while returning the new index
		deltaRow, deltaCol, newDirection, geneIndex, newR, newC := shark.UseGenomeToMove(currEco, i, j)
		isMoving := deltaRow != 0 || deltaCol != 0
		shark.DecreaseEnergy(geneIndex, isMoving)

		(*currEco)[i][j].predator = nil
		
		if shark.energy > 0 {
			currEco[newR][newC].predator.lastDirection = newDirection
			currEco[newR][newC].predator = shark		
		}
		
		// 2. FEEDING:
		// Check to eat fish or not
		shark.FeedShark(currEco, newR, newC)

		//3. AGE
		shark.UpdateAge() //Just add one

		//4. Reproduction
		if shark.CheckAge(ageThresholdPredator) && shark.CheckEnergy(energyThresholdPredator) {
			var babyShark Predator
			freeUnits := GetAvailableUnits(currEco, i, j)
			if len(freeUnits) != 0 {
				deltaX, deltaY := pickUnit(&freeUnits)
				(*currEco)[i+deltaX][j+deltaY] = &babyShark
				shark.Reproduce(&babyShark)
			}
			
		}

	}

}

// UseGenomeToMovePredator()
func (shark *Predator) UseGenomeToMove(currentEcosystem *Ecosystem, i, j int) (int, int, int, int, int, int) {
	var moveDeltas OrderedPair
	var geneIndex, newDirection, newI, newJ int
	isFreeUnitFlag := false
	numTries := 0

	// 20 is the threshold for max number of tries we get to reselect a gene for movement
	// if numberTries >= 20 and isFreeUnitFlag is still false
	// the prey doesn't move
	for !shark.isFreeUnit() && numTries < 20 {
		r := rand.Float64()
		geneIndex = 0
		runningSum := 0.0
		for idx, gene := range shark.genome {
			runningSum += float64(gene)
			if runningSum >= r {
				geneIndex = idx
				break
			}
		}
		newDirection := (shark.lastDirection + geneIndex) % 8
		moveDeltas = deltas[newDirection]
		numRows := len(*currentEcosystem)
		numCols := len((*currentEcosystem)[0])
		newI = GetIndex(i, moveDeltas.row, numRows)
		newJ = GetIndex(j, moveDeltas.col, numCols)


		//This check if the unit is free or not
		isFreeUnitFlag = isFreeUnit(currentEcosystem, i, j)
		numTries += 1
	}
	// if numTries >= 20 and still haven't find a free unit, we don't move
	if !shark.isFreeUnit() {
		geneIndex = 0
		newDirection = currentPredator.lastDirection
		moveDeltas.row, moveDeltas.col = 0, 0
	}

	//lastDirection will be updated with my new direction
	return moveDeltas.row, moveDeltas.col, newDirection, geneIndex, newI, newJ
}

func (shark *Predator) isFreeUnit(currEco *Ecosystem, i, j int) bool {
	return (*currEco)[i][j].predator == nil
}

func (shark *Predator) FeedShark(currEco *Ecosystem, x, y int) {
	if (*currEco)[x][y].prey != nil {
		(*currEco)[x][y].prey = nil
	}
	shark.IncreaseEngeryAfterMeal() //increase energy after eating a fish

}

func (shark *Predator) IncreaseEngeryAfterMeal() {
	shark.Organism.energy += 1
}

func GetAvailableUnits(currEco *Ecosystem, r, c int) []int {
	var units []int
	var n int
	for i := r - 1; i <= r+1; i++ {

		if i < 0 {
			i = len(*currEco) - 1
		}
		if i == len(*currEco) {
			i = 0
		}

		for j := c - 1; j <= c+1; j++ {
			if j < 0 {
				j = len(*currEco) - 1
			}
			if j == len(*currEco) {
				j = 0

				if IsItAvailable((*currEco)[i][j], true) {
					n = GetUnit(r, c, i, j, len(*currEco))
					units = append(units, n)
				}
			}
		}
	}
	return units
}


func (shark *Predator) Reproduce(babyShark *Predator) {
	//Already check age and energy!!!!
	shark.Organism.age = 0
	babyShark.Organism.energy = shark.Organism.energy / 2
	shark.Organism.energy /= 2
	babyShark.Organism.genome = shark.Organism.genome // Check if the array needs to be copied manually.
	UpdateDirection(&shark.Organism, &babyShark.Organism)
	UpdateGenome(&babyShark.Organism)

}

func (shark *Predator) CheckAge(threshold int) bool {
	return shark.Organism.age >= threshold
}

func (shark *Predator) CheckEnergy(threshold int) bool {
	return shark.Organism.energy >= threshold
}

func (shark *Predator) UpdateAge() {
	shark.Organism.age += 1
}

func (shark *Predator) DecreaseEnergy(geneIndex int, isMoving bool) {
	if isMoving {
		shark.energy -= energyCosts[geneIndex]
	}
}
