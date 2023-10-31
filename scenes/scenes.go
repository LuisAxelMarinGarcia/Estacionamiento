package scenes

import (
	"time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"carro/models"
	"carro/views"
)



func Run() {

	models.Init()

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Parking Lot Simulation",
		Bounds: pixel.R(0, 0, 800, 600),
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for car := range models.CarChannel {
			models.LaneMutex.Lock()
			for _, occupied := range models.LaneStatus {
				if !occupied {
					break
				}
			}
			models.LaneMutex.Unlock()

			go models.Lane(car.ID) 
		}
	}()
	
	
	for !win.Closed() {
		win.Clear(colornames.White)
		views.DrawParkingLot(win, models.GetCars())
		win.Update()
		models.CarsMutex.Lock()
		models.MoveCarsLogic()
		models.CarsMutex.Unlock()

		time.Sleep(16 * time.Millisecond)
	}
}


