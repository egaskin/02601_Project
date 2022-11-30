package main

import (
	"canvas"
	"fmt"
	"image"
)

// AnimateSystem takes a slice pointers to Ecosystem objects along with a canvas width
// parameter and a frequency parameter.
// Every frequency steps, it generates a slice of images corresponding to drawing each Ecosystems
// on a canvasWidth x canvasWidth canvas.
// A scaling factor is a final input that is used to scale the objects to be big enough to see
func AnimateSystem(allEcosystems []*Ecosystem, canvasWidth, frequency int, scalingFactor float64) []image.Image {
	images := make([]image.Image, 0)

	if len(allEcosystems) == 0 {
		panic("Error: no Ecosystems objects present in AnimateSystem.")
	}

	numberImages := len(allEcosystems)
	// for every universe, draw to canvas and grab the image
	for i := range allEcosystems {
		if i%frequency == 0 {
			// fmt.Println(i)
			images = append(images, allEcosystems[i].DrawToCanvas(canvasWidth, scalingFactor))
		}
		// print status of image drawing
		if (numberImages / 10) != 0 {
			if i%(numberImages/10) == 0 {
				fmt.Println("Drawing is", 100*float64(i)/float64(numberImages), "percent complete")
			}
		}
	}

	return images
}

// DrawToCanvas generates the image corresponding to a canvas after drawing a Ecosystem
// object's bodies on a square canvas that is canvasWidth pixels x canvasWidth pixels.
// A scaling factor is used to ensure objects are large enough
func (eco *Ecosystem) DrawToCanvas(canvasWidth int, scalingFactor float64) image.Image {
	if eco == nil {
		panic("Can't draw a nil Ecosystem.")
	}

	// get the size of the board
	numRows := eco.CountRows()
	numCols := eco.CountCols()
	unitWidth := 1

	// set a new square canvas
	// c := canvas.CreateNewCanvas(canvasWidth, canvasWidth)
	c := canvas.CreateNewCanvas(numRows, numCols)

	// create a black background
	c.SetFillColor(canvas.MakeColor(0, 0, 0))
	c.ClearRect(0, 0, canvasWidth, canvasWidth)
	c.Fill()

	// colors for each unit type
	var food_red uint8 = 255
	var food_green uint8 = 0
	var food_blue uint8 = 0
	foodColor := canvas.MakeColor(food_red, food_green, food_blue)

	// var prey_red uint8 = 0
	// var prey_green uint8 = 255
	// var prey_blue uint8 = 0
	// preyColor := canvas.MakeColor(prey_red, prey_green, prey_blue)

	// var prey_red uint8 = 0
	// var prey_green uint8 = 0
	// var prey_blue uint8 = 255
	// predColor := canvas.MakeColor(pred_red, pred_green, pred_blue)

	// range over all the Units and draw them.
	for i := range *eco {
		for j := range (*eco)[i] {
			curUnit := (*eco)[i][j]
			if curUnit.food.isPresent {
				// fmt.Println("this ran")
				c.SetFillColor(foodColor)

				// make the food smaller
				x := j * unitWidth
				y := i * unitWidth
				c.ClearRect(x, y, x+unitWidth, y+unitWidth)
				c.Fill()
			}
			// else if

			// c.SetFillColor(canvas.MakeColor(b.red, b.green, b.blue))
			// cx := (b.position.x / u.width) * float64(canvasWidth)
			// cy := (b.position.y / u.width) * float64(canvasWidth)
			// r := scalingFactor * (b.radius / u.width) * float64(canvasWidth)
			// c.Circle(cx, cy, r)
			// c.Fill()
		}

	}
	// we want to return an image!
	return c.GetImage()
}
