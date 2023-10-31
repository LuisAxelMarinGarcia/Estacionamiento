package scenes

import (
	"time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
	"carro/models"
	"carro/views"
)



func Run() {

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Parking Lot Simulation",
		Bounds: pixel.R(0, 0, 800, 600),
	})
	if err != nil {
		panic(err)
	}

	go func() {
		i := 1  
		for {   
			models.LaneMutex.Lock()
			for _, occupied := range models.LaneStatus {
				if !occupied {
					break
				}
			}
			models.LaneMutex.Unlock()
	
			go models.Lane(i)
			i++  
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500) + 500))
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


