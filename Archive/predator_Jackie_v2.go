package main

import "math/rand"

// UpdatePredator is a Predator method which will take a Predator input and update the position, initiate eating, reproduction, and age accordingly
func (shark *Predator) UpdatePredator(currEco, nextEco *Ecosystem, i, j int) {

	if shark.Organism.energy == 0 {

		(*nextEco)[i][j].predator = nil

	} else {

		//1. Update Position first if energy is allowed
		// This UpdatePosition will scan through all 7 units, give a list of available units, and use GENOME to update
		//	We prioritize the GENOME instead of the fish
		// This function will UpdatePredatorPosition while returning the new index
		x, y := shark.UpdatePredatorPosition(currEco, nextEco, i, j)

		// UseGenomeToMovePredator(currEco, (*nextEco)[i][j].predator,)

		// 2. FEEDING:
		// Check to eat fish or not
		shark.FeedShark(currEco, nextEco, x, y)

		//3. Reproduction

		if shark.CheckAge(20) && shark.CheckEnergy(1000) {
			var babyShark Predator
			shark.Reproduce(&babyShark)
		}

		//UpdateAge
		shark.UpdateAge() //Just add one
	}

}

// UseGenomeToMovePredator()
func UseGenomeToMovePredator(currentEcosystem *Ecosystem, currentPredator *Predator, i, j int) (int, int, int, int, int, int) {
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
		for idx, gene := range currentPredator.genome {
			runningSum += float64(gene)
			if runningSum >= r {
				geneIndex = idx
				break
			}
		}
		newDirection := (currentPredator.lastDirection + geneIndex) % 8
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
		newDirection = currentPredator.lastDirection
		moveDeltas.row, moveDeltas.col = 0, 0
	}

	//lastDirection will be updated with my new direction
	return moveDeltas.row, moveDeltas.col, newDirection, geneIndex, newI, newJ
}

// UpdatePredatorPosition is the Predator method which will update the position and energy of the Predator according to the food
// func (shark *Predator) UpdatePredatorPosition(currEco, nextEco *Ecosystem, r, c int) (int, int) {

// 	//1. Grab all available units around it
// 	vacantUnits := GetAvailableUnits(currEco, r, c)

// 	//2. Check the genome and give the best move
// 	tempGenome := make([]Gene, 8)
// 	for i := range shark.Organism.genome {
// 		tempGenome[i] = shark.Organism.genome[i]
// 	}
// 	sortedGenome := QuicksortGenome(tempGenome)

// 	sortedGenomeIndex := GiveSortedIndex(sortedGenome, tempGenome)

// 	var r0, c0 int

// 	if len(vacantUnits) != 0 {
// 		var nextMove int
// 		for i := range sortedGenomeIndex {
// 			if IsPresent(vacantUnits, sortedGenomeIndex[i]) {
// 				nextMove = sortedGenomeIndex[i]
// 				break
// 			}
// 		}

// 		r0, c0 = ConvertMovingIndices(nextMove, r, c)

// 		//After moving, decrease the energy
// 		shark.DecreaseEnergy()
// 	} else {
// 		r0 = r
// 		c0 = c
// 	}
// 	(*nextEco)[r0][c0].predator = shark

// 	return r0, c0
// }

func GiveSortedIndex(g1, g0 []Gene) []int {
	var index []int
	for i := range g1 {
		for j := range g0 {
			if g1[i] == g0[j] {
				index[i] = j
			}
		}
	}
	return index
}

func QuicksortGenome(genome []Gene) []Gene {

	var sortedGenome []Gene
	var left, right []Gene

	//randomly choose an pivot
	pivotIdx := PickPivot(len(genome))
	pivot := genome[pivotIdx]

	//Range over the genome, and compare with pivot
	for i := range genome {
		//Not include pivot
		if i != pivotIdx && genome[i] < pivot {
			left = append(left, genome[i])
		}
		if i != pivotIdx && genome[i] >= pivot {
			right = append(right, genome[i])
		}

	}
	left = QuicksortGenome(left)
	right = QuicksortGenome(right)

	sortedGenome = left
	sortedGenome = append(sortedGenome, pivot)
	sortedGenome = append(sortedGenome, right...)

	return sortedGenome
}

func PickPivot(len int) int {
	return rand.Intn(len)
}

func (shark *Predator) FeedShark(currEco, nextEco *Ecosystem, x, y int) {
	if (*currEco)[x][y].prey != nil {
		(*nextEco)[x][y].prey = nil
	}
	shark.IncreaseEngeryAfterMeal() //increase energy after eating a fish

}

func (shark *Predator) IncreaseEngeryAfterMeal() {
	shark.Organism.energy += 1
}

func IsPresent(vacantUnits []int, genomeIndex int) bool {
	for i := range vacantUnits {
		if genomeIndex == vacantUnits[i] {
			return true
		}
	}
	return false
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

func ConvertMovingIndices(nextMove, r, c int) (int, int) {
	var i, j int
	if nextMove == 0 {
		i = r - 1
		j = c - 1
	} else if nextMove == 1 {
		i = r - 1
		j = c
	} else if nextMove == 2 {
		i = r - 1
		j = c + 1
	} else if nextMove == 3 {
		i = r
		j = c - 1
	} else if nextMove == 4 {
		i = r
		j = c + 1
	} else if nextMove == 5 {
		i = r + 1
		j = c - 1
	} else if nextMove == 6 {
		i = r + 1
		j = c
	} else if nextMove == 7 {
		i = r + 1
		j = c + 1
	}
	return i, j
}

func GetUnit(r, c, i, j, n int) int {
	var unit int
	rowDelta := r - i
	colDelta := c - j

	//edge case
	if rowDelta < -1 { //first row
		rowDelta = 1
	}
	if rowDelta > 1 { //last row
		rowDelta = -1
	}
	if colDelta < -1 { //first col
		colDelta = 1
	}
	if colDelta > 1 { //last col
		colDelta = -1
	}

	if rowDelta == 1 && colDelta == 1 {
		unit = 0
	} else if rowDelta == 1 && colDelta == 0 {
		unit = 1
	} else if rowDelta == 1 && colDelta == -1 {
		unit = 2
	} else if rowDelta == 0 && colDelta == 1 {
		unit = 3
	} else if rowDelta == 0 && colDelta == -1 {
		unit = 4
	} else if rowDelta == -1 && colDelta == 1 {
		unit = 5
	} else if rowDelta == -1 && colDelta == 0 {
		unit = 6
	} else if rowDelta == -1 && colDelta == -1 {
		unit = 7
	}
	return unit
}

func IsItAvailable(unit *Unit, IsThisAPredator bool) bool {
	//Check if there is any predator
	return unit.predator == nil
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

func (shark *Predator) DecreaseEnergy() {
	shark.Organism.energy -= 1
}
