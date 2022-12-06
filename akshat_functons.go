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

func ReproducePrey(p *Prey) *Prey {
	//This function will only be called if the age and energy and requirements are met. Check these requirements before calling this function.
	var child Prey
	p.Organism.age = 0
	child.Organism.energy = p.Organism.energy / 2
	p.Organism.energy /= 2
	child.Organism.genome = p.Organism.genome // Check if the array needs to be copied manually.
	UpdateDirection(&p.Organism, &child.Organism)
	UpdateGenome(&child.Organism)
	return &child
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

	currentPrey := (*currentEcosystem)[i][j].prey
	// only move the prey if hasn't moved this generation
	if currentPrey.lastGenUpdated != currGen {
		currentPrey.lastGenUpdated = currGen
		UpdateAgePrey(currentPrey)
		if (*currentEcosystem)[i][j].prey.energy >= energyThreshold && (*currentEcosystem)[i][j].prey.age >= ageThreshold {
			freeUnits := GetAvailableUnits(currentEcosystem, i, j)
			if len(freeUnits) != 0 {
				deltaX, deltaY := pickUnit(&freeUnits)
				(*currentEcosystem)[i+deltaX][j+deltaY].prey = ReproducePrey(currentPrey)
			}

		}
		MovePrey(currentEcosystem, i, j)
	}
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
