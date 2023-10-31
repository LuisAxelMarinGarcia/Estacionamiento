	package views

	import (
		"github.com/faiface/pixel"
		"github.com/faiface/pixel/imdraw"
		"github.com/faiface/pixel/pixelgl"
		"golang.org/x/image/colornames"
		"carro/models"
		"image"
		_ "image/png"
		"os"
	)

	var (
		background *pixel.Sprite
		bgPicture  pixel.Picture
	)
	
	func loadBackground() {
		file, err := os.Open("Assets/background.png")
		if err != nil {
			panic(err)
		}
		defer file.Close()
	
		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}
	
		bgPicture = pixel.PictureDataFromImage(img)
		background = pixel.NewSprite(bgPicture, bgPicture.Bounds())
	}
	
	func DrawParkingLot(win *pixelgl.Window, cars []models.Car) {
		if background == nil {
			loadBackground()
		}
	
		background.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		imd := imdraw.New(nil)
		imd.Color = colornames.White

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
	
			topLeft := car.Position.Add(pixel.V(-carWidth, -carHeight))
			topRight := car.Position.Add(pixel.V(carWidth, -carHeight))
			bottomLeft := car.Position.Add(pixel.V(-carWidth, carHeight))
			bottomRight := car.Position.Add(pixel.V(carWidth, carHeight))
	
			imd.Push(topLeft, topRight, bottomRight, bottomLeft)
			imd.Polygon(0)
		}
	
		imd.Draw(win)
	}