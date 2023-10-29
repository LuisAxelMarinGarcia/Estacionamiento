package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func drawParkingLot(win *pixelgl.Window, cars []Car) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black

	// Dibuja el contorno del estacionamiento
	imd.Push(pixel.V(100, 500), pixel.V(700, 500))
	imd.Line(2)
	imd.Push(pixel.V(100, 100), pixel.V(700, 100))
	imd.Line(2)

	// Lado derecho del contorno
	imd.Push(pixel.V(700, 100), pixel.V(700, 500))
	imd.Line(2)

	// Dibuja los carriles de estacionamiento
	for i := 0.0; i < 4.0; i++ {
		xOffset := 100.0 + i*(laneWidth)
		// Carriles superiores
		imd.Push(pixel.V(xOffset, 500), pixel.V(xOffset+laneWidth, 350))
		imd.Rectangle(2)
		// Carriles inferiores
		imd.Push(pixel.V(xOffset, 250), pixel.V(xOffset+laneWidth, 100))
		imd.Rectangle(2)
	}

    for _, car := range cars {
        imd.Color = colornames.Red
        // Asegúrate de que el carro se dibuje centrado en su posición
        imd.Push(car.Position.Add(pixel.V(-25, -25)), car.Position.Add(pixel.V(25, 25)))  // Asume que cada carro es un cuadrado de 50x50 pixels
        imd.Rectangle(2)
    }

    imd.Draw(win)
}