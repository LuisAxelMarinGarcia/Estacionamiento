	package views

	import (
		"github.com/faiface/pixel"
		"github.com/faiface/pixel/imdraw"
		"github.com/faiface/pixel/pixelgl"
		"golang.org/x/image/colornames"
		"carro/models"
	)

	func DrawParkingLot(win *pixelgl.Window, cars []models.Car) {
		imd := imdraw.New(nil)
		imd.Color = colornames.Black

		imd.Push(pixel.V(100, 500), pixel.V(700, 500))
		imd.Line(2)
		imd.Push(pixel.V(100, 100), pixel.V(700, 100))
		imd.Line(2)

		imd.Push(pixel.V(700, 100), pixel.V(700, 500))
		imd.Line(2)

		parkingWidth := 600.0 
		laneWidth := parkingWidth / 10 

		for i := 0.0; i < 10.0; i++ {
			xOffset := 100.0 + i*laneWidth 
			imd.Push(pixel.V(xOffset, 500), pixel.V(xOffset, 350))
			imd.Line(2)  

			imd.Push(pixel.V(xOffset, 250), pixel.V(xOffset, 100))
			imd.Line(2)  
		}

		carWidth := laneWidth / 4  
		carHeight := laneWidth / 4 

		for _, car := range cars {
			imd.Color = colornames.Red
			imd.Push(car.Position.Add(pixel.V(-carWidth, -carHeight)), car.Position.Add(pixel.V(carWidth, carHeight))) 
			imd.Rectangle(2)
		}

		imd.Draw(win)
	}
