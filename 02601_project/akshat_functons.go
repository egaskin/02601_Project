package main

import "math/rand"

//UpdateAge takes in a pointer to a prey object, updates it age by incrementing it by one and then returns it.
func UpdateAgePrey(p *Prey) {
	p.Organism.age += 1
}

//UpdateAge takes in a pointer to a prey object, updates it age by incrementing it by one and then returns it.
func UpdateAgePredator(p *Predator) {
	p.Organism.age += 1
}

func ReproducePrey(parent, child *Prey) {
	//This function will only be called if the age and energy and requirements are met. Check these requirements before calling this function.
	parent.Organism.age = 0
	child.Organism.energy = parent.Organism.energy / 2
	parent.Organism.energy /= 2
	child.Organism.genome = parent.Organism.genome // Check if the array needs to be copied manually.
	UpdateDirection(&parent.Organism, &child.Organism)
	UpdateGenome(&child.Organism)
}

//UpdateDirection updates the direction of that the child is moving in based on the parents genome and direction of movement
func UpdateDirection(parent, child *Organism) {
	r := rand.Float64()
	var sum Gene
	index := 0
	for i := range parent.genome {
		if Gene(r) < sum {
			index = i
			break
		} else {
			sum += parent.genome[i]
		}
	}
	child.lastDirection = (parent.lastDirection + index) % 8
}

//UpdateGenome updates the genome of the child based on the last known movement.
func UpdateGenome(currentOrganism *Organism) {
	currentDirection := currentOrganism.lastDirection
	delta := 0.1
	currentOrganism.genome[currentDirection] += Gene(delta) * currentOrganism.genome[currentDirection]
	for i := range currentOrganism.genome {
		if i != currentDirection {
			if currentOrganism.genome[i]-Gene(delta)*currentOrganism.genome[currentDirection] > 0 {
				currentOrganism.genome[i] -= Gene(delta) * currentOrganism.genome[currentDirection] / 7.0
			}
		}
	}
}

func ReproducePredator(p *Predator) *Predator {
	//This function will only be called if the age and energy and requirements are met. Check these requirements before calling this function.
	var child Predator
	p.Organism.age = 0
	child.Organism.energy = p.Organism.energy / 2
	p.Organism.energy /= 2
	child.Organism.genome = p.Organism.genome // Check if the array needs to be copied manually.
	UpdateDirection(&p.Organism, &child.Organism)
	UpdateGenome(&child.Organism)
	return &child
}

func UpdatePrey(currentEcosystem *Ecosystem, i, j, currGen int) {
	numRows := len(*currentEcosystem)
	numCols := len((*currentEcosystem)[0])
	currentPrey := (*currentEcosystem)[i][j].prey
	// note we have moved the prey this timestep/generation
	currentPrey.lastGenUpdated = currGen

	if currentPrey.Organism.energy <= 0 {
		(*currentEcosystem)[i][j].prey = nil
		return
	}

	UpdateAgePrey(currentPrey)
	if (*currentEcosystem)[i][j].prey.energy >= energyThresholdPrey && (*currentEcosystem)[i][j].prey.age >= ageThresholdPrey {
		var babyPrey Prey
		freeUnits := GetAvailableUnits(currentEcosystem, i, j)
		if len(freeUnits) != 0 {
			deltaX, deltaY := pickUnit(&freeUnits)
			newI := GetIndex(i, deltaX, numRows)
			newJ := GetIndex(j, deltaY, numCols)
			(*currentEcosystem)[newI][newJ].prey = &babyPrey
			ReproducePrey(currentPrey, &babyPrey)
		}

	}
	MovePrey(currentEcosystem, i, j)

}

func pickUnit(freeUnits *[]int) (r, c int) {
	length := len(*freeUnits)
	random := rand.Intn(length)
	chosenUnit := (*freeUnits)[random]
	return GetIndices(&chosenUnit)
}

func GetIndices(chosenUnit *int) (r, c int) {
	if *chosenUnit == 0 {
		return -1, -1
	}
	if *chosenUnit == 1 {
		return -1, 0
	}
	if *chosenUnit == 2 {
		return -1, 1
	}
	if *chosenUnit == 3 {
		return 0, 1
	}
	if *chosenUnit == 4 {
		return 1, 1
	}
	if *chosenUnit == 5 {
		return 1, 0
	}
	if *chosenUnit == 6 {
		return 1, -1
	}
	if *chosenUnit == 7 {
		return 0, -1
	}
	return 0, 0
}

// Akshat: InitializePreyAndPredator
// Randomly generate numPrey and numPred predators in the initialEcosystem.
// Functions written by Akshat
func InitializePreyAndPredator(numRows, numCols, numPrey, numPred int, newEco *Ecosystem) {
	// Akshat wrote these: Randomly initialize the prey and predators
	count_Prey := 0
	count_Pred := 0

	for count_Pred < numPred {
		i := rand.Intn(numRows)
		j := rand.Intn(numCols)
		if (*newEco)[i][j].predator == nil {
			(*newEco)[i][j].predator = CreatePredator()
		}
		count_Pred += 1
	}

	for count_Prey < numPrey {
		i := rand.Intn(numRows)
		j := rand.Intn(numCols)
		if (*newEco)[i][j].prey == nil && (*newEco)[i][j].predator == nil {
			(*newEco)[i][j].prey = CreatePrey()
			count_Prey += 1
		}
	}
}

// Akshat: CreatePrey initializes the Prey object
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

// Akshat: CreatePrey initializes the Predator object
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

// Akshat: CreateGenome creates the first version of the genome and returns it.
func CreateGenome() [8]Gene {
	var newGenome [8]Gene
	for i := range newGenome {
		newGenome[i] = Gene(0.125)
	}
	return newGenome
}
