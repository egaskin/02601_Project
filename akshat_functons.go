package main

import (
	"fmt"
	"gifhelper"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {

	// var numRows int = 250
	// var numCols int = 250
	var numRows int = 100
	var numCols int = 100
	var numPrey int = 0
	var numPred int = 0
	var initialEcosystem Ecosystem = InitializeEcosystem(numRows, numCols, numPrey, numPred)
	var totalTimesteps int = 100
	var foodRule string = "gardenOfEden"

	// seed the PRNG approximately randomly
	rand.Seed(time.Now().UnixNano())
	allEcosystems := SimulateEcosystemEvolution(&initialEcosystem, totalTimesteps, foodRule)

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
	fmt.Println("Also need to initialize the food randomly!!!")
	for i := 0; i < numRows; i++ {
		newEco[i] = make([]*Unit, numCols)
		for j := 0; j < numCols; j++ {

			// initialize the pointer
			newEco[i][j] = new(Unit)
		}
	}
	//Randomly initialize the prey and predators
	count_Prey := 0
	count_Pred := 0

	for count_Pred <= numPred {
		i := rand.Intn(numRows + 1)
		j := rand.Intn(numCols + 1)
		if newEco[i][j].prey == nil {
			newEco[i][j].predator = CreatePredator()
		}
		count_Pred += 1
	}

	for count_Prey <= numPrey {
		i := rand.Intn(numRows + 1)
		j := rand.Intn(numCols + 1)
		if newEco[i][j].prey == nil && newEco[i][j].predator == nil {
			newEco[i][j].prey = CreatePrey()
			count_Prey += 1
		}
	}

	return Ecosystem(newEco)
}

// CreatePrey initializes the Prey object
func CreatePrey() *Prey {
	var newPrey Prey
	newPrey.Organism.age = 0
	newPrey.Organism.energy = 5
	newPrey.Organism.age = 0
	newPrey.Organism.genome = CreateGenome()
	newPrey.Organism.lastGenUpdated = 0
	newPrey.Organism.lastDirection = 0
	return &newPrey
}

// CreatePrey initializes the Predator object
func CreatePredator() *Predator {
	var newPredator Predator
	newPredator.Organism.age = 0
	newPredator.Organism.energy = 5
	newPredator.Organism.age = 0
	newPredator.Organism.genome = CreateGenome()
	newPredator.Organism.lastGenUpdated = 0
	newPredator.Organism.lastDirection = 0
	return &newPredator

}

// CreateGenome creates the first version of the genome and returns it.
func CreateGenome() [8]Gene {
	var newGenome [8]Gene
	for i := range newGenome {
		newGenome[i] = Gene(0.125)
	}
	return newGenome
}

// SimulateEcosystemEvolution() takes an Ecosystem pointer initialEcosystem, int totalTimesteps, and string foodRule. The function sequentially simulates initialEcosystem evolving over the course of totalTimesteps generations and saves each Ecosystem into a collection for the output, a slice of Ecosystem pointers called allEcosystems.
func SimulateEcosystemEvolution(initialEcosystem *Ecosystem, totalTimesteps int, foodRule string) []*Ecosystem {

	// initialize the slice to contain all the Ecosystems
	allEcosystems := make([]*Ecosystem, totalTimesteps+1)
	allEcosystems[0] = initialEcosystem
	// keep track of start of simulation
	var start time.Time = time.Now()

	// sequentially update allEcosystems (serially) and save them
	for i := 1; i <= totalTimesteps; i++ {
		allEcosystems[i] = UpdateEcosystem(allEcosystems[i-1], foodRule, i)

		// print status of simulation
		if (totalTimesteps / 10) != 0 {
			if i%(totalTimesteps/10) == 0 || i == 1 {
				fmt.Println("Simulation is", 100*float64(i)/float64(totalTimesteps), "percent complete")
				elapsed := time.Since(start)
				log.Printf("This took total %s\n\n", elapsed)
			}
		}
	}

	return allEcosystems
}

func UpdateEcosystem(prevEcosystem *Ecosystem, foodRule string, curGen int) *Ecosystem {

	// initialize the nextEcosystem
	var nextEcosystem *Ecosystem = DeepCopyEcosystem(prevEcosystem)

	// get the numRows and numCols of the ecosystem
	numRows := nextEcosystem.CountRows()
	numCols := nextEcosystem.CountCols()
	arrayOfIndices := MakeIndicesArray(nextEcosystem)

	// get ready to loop through all the possible indices of the ecosystem using a while loop
	totalUpdates := numCols * numRows
	k := 0
	for k < totalUpdates {

		i, j := ChooseRandomIndices(numRows, numCols)

		// only update if the chosenRow, chosenCol is a valid location in the ecosystem
		if i != -1 && j != -1 {
			// only update i when we visit a valid location
			k++
			// update the indices array
			UpdateIndices(&arrayOfIndices, i, j)

			// get a pointer to the current Unit that we need new values for (nextEcosystem)
			currentUnit := (*nextEcosystem)[i][j]

			// this variable keeps track of whether there is more than one thing in the current Unit or not. there should only be either prey or predator or food. never multiple things in the Unit at once
			thereIsPred := false

			// Update the currentUnit based on the nextEcosystem! since we want the system to change as prey and food or disappearing (so each prey/predator is competing to get to their respective food source first)

			// only perform this operation is Unit contains predator
			if (*currentUnit).predator != nil {

				// skip already updated predator.
				if (*currentUnit).predator.lastGenUpdated != curGen {
					// currentUnit.predator.UpdatePred(nextEcosystem, i, j)
					currentUnit.predator.lastGenUpdated = curGen
				}

				thereIsPred = true

			}
			if (*currentUnit).prey != nil { // only perform this operation if Unit contains predator

				// skip already updated prey.
				if (*currentUnit).prey.lastGenUpdated != curGen {
					//currentUnit.prey.UpdatePrey(nextEcosystem, i, j)
					currentUnit.prey.lastGenUpdated = curGen
				}

				if thereIsPred {
					panicStatement := "there is a predator and prey in the Unit row, col " + strconv.Itoa(i) + "," + strconv.Itoa(j)
					panic(panicStatement)
				}

			}

			// we allow predator and prey stacking on top of food
			if !(*currentUnit).food.isPresent { // if we made it here that means there can only be food in the current Unit. skip if the food is already true. note: we don't need to check the lastGenUpdated because food will be false if this is ran.

				// determine whether food appears randomly for the prey. GeneratePreyFoodRandomly() will update both fields of the food, if food is generated. otherwise it will leave it false.
				currentUnit.GeneratePreyFoodProbabilistically(foodRule, i, j, nextEcosystem)
				// currentUnit.food.lastGenUpdated = curGen

			}

		}
	}

	return nextEcosystem
}

// Input: the number of rows and number of cols to choose from, numRows and numCols
// Output: two integers, randomly choosen between for intervals [0,numRows) and [0,numCols)
func ChooseRandomIndices(numRows, numCols int) (int, int) {
	chosenRow := rand.Intn(numRows)
	chosenCol := rand.Intn(numCols)

	return chosenRow, chosenCol
}

// Input: an Ecosystem pointer
// Output: a 2D array containing every combination of row and col (all the indices of that Ecosystem)
func MakeIndicesArray(someEcosystem *Ecosystem) [][]OrderedPair {

	// get the number of rows and cols in the ecosystem
	numRows := someEcosystem.CountRows()
	numCols := someEcosystem.CountCols()

	// initialize the arrays to contain all the rows and column indices in the ecosystem
	arrayOfIndices := make([][]OrderedPair, numRows)

	// range over all the arrayOfIndices and set equal to the index
	for i := range arrayOfIndices {
		arrayOfIndices[i] = make([]OrderedPair, numCols)
		for j := range arrayOfIndices[i] {
			arrayOfIndices[i][j].row = i
			arrayOfIndices[i][j].col = j
		}
	}
	return arrayOfIndices
}

// Input: indicesAvailable, a 2D array of OrderedPair containing the indices available for the ecosystem to choose. dont confuse the indices of this array, with the indices available!
// Output: none! operates on a pointer to makes updated indicesAvailable
func UpdateIndices(indicesAvailable *[][]OrderedPair, rowIndex, colIndex int) {
	// set the index values at that location to -1, -1 so that we know we've visited that location already
	(*indicesAvailable)[rowIndex][colIndex].row = -1
	(*indicesAvailable)[rowIndex][colIndex].col = -1
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

