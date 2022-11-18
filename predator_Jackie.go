package main

// UpdatePredator is a Predator method which will take a Predator input and update the position, initiate eating, reproduction, and age accordingly
func (shark *Predator) UpdatePredator(nextEco Ecosystem, i, j int) {
	
	
	shark.UpdatePredatorPosition(i, j) {
		
		//Scan for the fish
		var edibleFishes []*Prey
		nemo := FindingNemo(i, j) { //return slice of Orderedpair Index of an edible Prey location around the Predator
			//Will take into consideration of genome later 
			//Scan through adjacent units of the currShark
			//Check if there is any fish
			//If there is fishes, let's check if we can eat the fish
			CheckIfEdible(fish) {
				bool
				if true //append to edibleFishes
				if false // move on to next one 
			}
			if CheckIfEdible == true { //append to edibleFish
			}
		}
		
		//If there is fishes, randomly choose one fish out of edibleFishes to eat		
		EatTheFish(edibleFishes) {
			index := rand.Int(len(edibleFishes))
			//nemo = index of edibleFishes[index]
			DeleteTheFish(index)
			shark.IncreaseEngeryAfterMeal() //increase energy after eating a fish
		}

		//If not fish
		var emptySpot []OrderPaired
			if len(edibleFishes) == 0 {
				//Find available position if there is no fish
				//Randomly choose t
			}
	

		shark.MovePredator(nemo) //Move Predator to the new location

	}
		
	// //I think EatFood should be inside UpdatePosition so I can access the info of the fish *********
	// shark.EatFood(i, j) //aka shark.FeedTime
	
	//Reproduction
	shark.Reproduce(i, j) {
		shark.CheckAge(threshold)
		shark.CheckEnergy(threshold)
		//If okay, Reproduce
		shark.AfterBirthEnergy()
		babyShark.AfterBirthEnergy()
	}
	
	//UpdateAge
	shark.UpdateAge() //Just add one
}

// UpdatePredatorPosition is the Predator method which will update the position and energy of the Predator according to the food
func (shark *Predator) UpdatePredatorPosition(r, c int) {
	//Range through the board and check unit around if there is any fish
	//If there is fish, then move toward that unit
	//If not, randomly move

	//Find the position (r and c) of the fish
	//If there is fish,return the fish's c and r
	//If not, return
	nemo := FindingNemo(r, c) //return Ordered Pair row and column of the fish
	//What if we have more than one fish

	//MovePredator takes the position of the next move and move the Predator there
	shark.MovePredator(nemo) // Move the Predator

	//After moving, decrease the energy
	shark.DecreaseEnergy()

}

func (shark *Predator) Reproduce(c, r, threshold int) *Predator {
	var babyShark Predator

	shark.CheckAge(threshold)
	shark.CheckEnergy(threshold)

	shark.AfterBirthEnergy()
	babyShark.AfterBirthEnergy()

	return &babyShark
}

func (shark *Predator) EatFood(i, j int) {

}

func (shark Predator) UpdateAge() {
	shark.age += 1
}
