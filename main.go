package main

import (
	"fmt"
	"gifhelper"
	"log"
	"strconv"
	"time"
)

func main() {

	var numRows int = 100
	var numCols int = 100
	var numPrey int = 0
	var numPred int = 0
	var initialEcosystem Ecosystem = InitializeEcosystem(numRows, numCols, numPrey, numPred)
	var numGens int = 1000
	var foodRule string = "eden"

	allEcosystems := SimulateEcosystemEvolution(&initialEcosystem, numGens, foodRule)

	canvasWidth := 1000
	frequency := 1
	scalingFactor := 1.0
	imageList := AnimateSystem(allEcosystems, canvasWidth, frequency, scalingFactor)

	gifhelper.ImagesToGIF(imageList, "ecosystem")
	fmt.Println("GIF drawn.")

}

func InitializeEcosystem(numRows, numCols, numPrey, numPred int) Ecosystem {

	// initialize newEco, which has numRows rows. the outer dimension
	newEco := make([][]*Unit, numRows)

	fmt.Println("NEED TO INITIALIZE PREY AND PREDATOR RANDOMLY WITHIN THIS. also lots to update in function.go like other food rules. and methods for prey and predator.")
	for i := 0; i < numRows; i++ {
		newEco[i] = make([]*Unit, numCols)
		for j := 0; j < numCols; j++ {

			// initialize the pointer
			newEco[i][j] = new(Unit)
		}
	}

	return newEco
}

// SimulateEcosystemEvolution() takes an Ecosystem pointer initialEcosystem, int numGens, and string foodRule. The function sequentially simulates initialEcosystem evolving over the course of numGens generations and saves each Ecosystem into a collection for the output, a slice of Ecosystem pointers called allEcosystems.
func SimulateEcosystemEvolution(initialEcosystem *Ecosystem, numGens int, foodRule string) []*Ecosystem {

	// initialize the slice to contain all the Ecosystems
	allEcosystems := make([]*Ecosystem, numGens+1)
	allEcosystems[0] = initialEcosystem
	// keep track of start of simulation
	var start time.Time = time.Now()

	// sequentially update allEcosystems (serially) and save them
	for i := 1; i <= numGens; i++ {
		allEcosystems[i] = UpdateEcosystem(allEcosystems[i-1], foodRule, i)

		// print status of simulation
		if (numGens / 10) != 0 {
			if i%(numGens/10) == 0 || i == 1 {
				fmt.Println("Simulation is", 100*float64(i)/float64(numGens), "percent complete")
				elapsed := time.Since(start)
				log.Printf("This took total %s\n\n", elapsed)
			}
		}
	}

	return allEcosystems
}

func UpdateEcosystem(prevEcosystem *Ecosystem, foodRule string, curGen int) *Ecosystem {
	// var row int = len(*prevEcosystem)
	// var col int = len((*prevEcosystem)[0])
	var nextEcosystem *Ecosystem = DeepCopyEcosystem(prevEcosystem)

	// range over all the units in the Ecosystem
	for i := range *nextEcosystem {
		for j := range (*nextEcosystem)[i] {
			// get a pointer to the current Unit that we need new values for (nextEcosystem)
			currentUnit := (*nextEcosystem)[i][j]

			// this variable keeps track of whether there is more than one thing in the current Unit or not. there should only be either prey or predator or food. never multiple things in the Unit at once
			oneThingRan := false

			// if (*currentUnit).prey != nil && (*currentUnit).predator != nil {
			// 	formatString := "(" + strconv.Itoa(i) + "," + strconv.Itoa(j) + ")"
			// 	panic("there is a prey and predator in Unit (i,j) = " + formatString)
			// } else if (*currentUnit).prey != nil && (*currentUnit).food.isPresent {
			// 	formatString := "(" + strconv.Itoa(i) + "," + strconv.Itoa(j) + ")"
			// 	panic("there is a prey and food in Unit (i,j) = " + formatString)
			// } else if (*currentUnit).predator != nil && (*currentUnit).food.isPresent {
			// 	formatString := "(" + strconv.Itoa(i) + "," + strconv.Itoa(j) + ")"
			// 	panic("there is a predator and food in Unit (i,j) = " + formatString)
			// }

			// Update the currentUnit based on the nextEcosystem! since we want the system to change as prey and food or disappearing (so each prey/predator is competing to get to their respective food source first)

			// only perform this operation is Unit contains predator
			if (*currentUnit).predator != nil {

				// skip already updated predator.
				if (*currentUnit).predator.lastGenUpdated != curGen {
					currentUnit.predator.UpdatePred(nextEcosystem, i, j)
					currentUnit.predator.lastGenUpdated = curGen
				}

				oneThingRan = true

			} else if (*currentUnit).prey != nil { // only perform this operation if Unit contains predator

				// skip already updated prey.
				if (*currentUnit).prey.lastGenUpdated != curGen {
					// currentUnit.prey.UpdatePrey(nextEcosystem, i, j)
					currentUnit.prey.lastGenUpdated = curGen
				}

				if oneThingRan {
					panicStatement := "two things ran for a single unit with row, col " + strconv.Itoa(i) + "," + strconv.Itoa(j)
					panic(panicStatement)
				}

				oneThingRan = true

			} else if !(*currentUnit).food.isPresent { // if we made it here that means there can only be food in the current Unit. skip if the food is already true. note: we don't need to check the lastGenUpdated because food will be false if this is ran.

				// determine whether food appears randomly for the prey. GeneratePreyFoodRandomly() will update both fields of the food, if food is generated. otherwise it will leave it false.
				currentUnit.GeneratePreyFoodProbabilistically(foodRule, i, j, nextEcosystem)
				// currentUnit.food.lastGenUpdated = curGen

				if oneThingRan {
					panicStatement := "two things ran for a single unit with row, col " + strconv.Itoa(i) + "," + strconv.Itoa(j)
					panic(panicStatement)
				}

			}

		}
	}

	return nextEcosystem
}

// assumes rectangular Ecosystem (rows all the same length)
func (someEcosystem *Ecosystem) CountRows() int {
	return len(*someEcosystem)
}

// assumes rectangular Ecosystem (cols all the same length)
func (someEcosystem *Ecosystem) CountCols() int {
	return len((*someEcosystem)[0])
}

func DeepCopyEcosystem(someEcosystem *Ecosystem) *Ecosystem {
	numCols := someEcosystem.CountCols()
	numRows := someEcosystem.CountRows()

	// make the rows of the copy (outermost dimension)
	var copyEcosystem Ecosystem = make([][]*Unit, numRows)

	// range over all the Units in someEcosystem's (every combo of row and col)
	for i := 0; i < numRows; i++ {
		// make the columns of the copy
		copyEcosystem[i] = make([]*Unit, numCols)
		for j := 0; j < numCols; j++ {
			// initialize the Unit that will go at i, j
			copyEcosystem[i][j] = new(Unit)

			// copy the corresponding fields of the Unit (deep copy)
			copyEcosystem[i][j].food = (*someEcosystem)[i][j].food

			// only attempt to copy if its there
			if (*someEcosystem)[i][j].prey != nil {
				copyEcosystem[i][j].prey = (*someEcosystem)[i][j].prey.DeepCopyOrganism()
			}

			// only attempt to copy if its there
			if (*someEcosystem)[i][j].predator != nil {
				copyEcosystem[i][j].predator = (*someEcosystem)[i][j].predator.DeepCopyOrganism()
			}
		}
	}

	return &copyEcosystem
}

func (somePrey *Prey) DeepCopyOrganism() *Prey {
	var preyCopy Prey
	fmt.Println("THIS IS:", somePrey)
	preyCopy.age = somePrey.age
	preyCopy.energy = somePrey.energy

	// range over the genome and copy all its genes
	var copyGenome [8]Gene
	for i := range copyGenome {
		copyGenome[i] = (*somePrey).genome[i]
	}

	// assign it into preyCopy
	preyCopy.genome = copyGenome

	return &preyCopy // return a pointer to the new copy
}

func (somePred *Predator) DeepCopyOrganism() *Predator {
	var predCopy Predator

	predCopy.age = somePred.age
	predCopy.energy = somePred.energy

	// range over the genome and copy all its genes
	var copyGenome [8]Gene
	for i := range copyGenome {
		copyGenome[i] = (*somePred).genome[i]
	}

	// assign it into predCopy
	predCopy.genome = copyGenome

	return &predCopy // return a pointer to the new copy
}
