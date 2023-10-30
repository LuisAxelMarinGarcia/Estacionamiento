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

var allParked bool  // false si alguno no está estacionado, true si todos están estacionados
var exitIndex int


func Run() {
	models.Entrance = make(chan bool, 1) // canal con buffer de 1
	exitIndex = 0

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Parking Lot Simulation",
		Bounds: pixel.R(0, 0, 800, 600),
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for i := 1; i <= 8; i++ {
			models.LaneMutex.Lock()
			allLanesOccupied := true
			for _, occupied := range models.LaneStatus {
				if !occupied {
					allLanesOccupied = false
					break
				}
			}
			models.LaneMutex.Unlock()
			
			if allLanesOccupied {
				break  // Detiene el spawn de nuevos carros si todos los carriles están ocupados
			}
	
			go models.Carr(i)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000) + 1000))  // Espera más tiempo antes de lanzar el próximo carro
		}
	}()
	
	
	for !win.Closed() {
		win.Clear(colornames.White)
		views.DrawParkingLot(win, models.GetCars())
		win.Update()
	
		models.CarsMutex.Lock()
		for i := len(models.Cars) - 1; i >= 0; i-- {
			for j := range models.Cars {
				if i != j {  // Asegúrate de no comparar el carro consigo mismo
					if models.CheckCollision(models.Cars[i].Position, models.Cars[j].Position) {
						continue  // Si hay colisión, salta al siguiente ciclo del bucle
					}
				}
			}
	
			if models.Cars[i].Position.X < 100 && models.Cars[i].Lane == -1 {
				models.Cars[i].Position.X += 2  // Mueve más lentamente hacia la derecha
				if models.Cars[i].Position.X > 100 {  // Asegura que el carro se detenga en X=100
					models.Cars[i].Position.X = 100
				}
			} else if models.Cars[i].Position.X == 100 && models.Cars[i].Lane == -1 {
				// La lógica de entrada existente...
			} else if models.Cars[i].Lane != -1 && !models.Cars[i].Parked {
				var targetX, targetY float64
				if models.Cars[i].Lane < 4 {
					// Carriles superiores
					targetX = 100.0 + float64(models.Cars[i].Lane)*views.LaneWidth + views.LaneWidth/2
					targetY = 350 + (500-350)/2  // Centro vertical de los carriles superiores
				} else {
					// Carriles inferiores
					targetX = 100.0 + float64(models.Cars[i].Lane-4)*views.LaneWidth + views.LaneWidth/2
					targetY = 100 + (250-100)/2  // Centro vertical de los carriles inferiores
				}
				if models.Cars[i].Position.X < targetX - 2 {
					models.Cars[i].Position.X += 2  // Mueve hacia la derecha
				} else if models.Cars[i].Position.Y < targetY - 2 && (targetX - models.Cars[i].Position.X) <= 2 {
					models.Cars[i].Position.Y += 2  // Mueve hacia arriba
				} else if models.Cars[i].Position.Y > targetY + 2 && (targetX - models.Cars[i].Position.X) <= 2 {
					models.Cars[i].Position.Y -= 2  // Mueve hacia abajo
				} else if (targetX - models.Cars[i].Position.X) <= 2 && (targetY - models.Cars[i].Position.Y) <= 2 {
					models.Cars[i].Position.X = targetX  // Ajusta la posición X a la posición objetivo
					models.Cars[i].Position.Y = targetY  // Ajusta la posición Y a la posición objetivo
					models.Cars[i].Parked = true  // Marca el carro como estacionado
					models.SetExitTime(&models.Cars[i])  // Establece el tiempo de salida del carro

					allParked = true  // asume que todos están estacionados
					for _, car := range models.Cars {
						if !car.Parked {
							allParked = false  // Si encuentra un carro que no está estacionado, establece allParked a false
							break
						}
					}
				}
				} else if models.Cars[i].Parked && time.Now().After(models.Cars[i].ExitTime) && allParked && i == exitIndex {
					if models.Cars[i].Lane < 4 {
						// Para carriles superiores
						if models.Cars[i].Position.Y > 300 {
							models.Cars[i].Position.Y -= 2  // Mueve hacia abajo
						} else if models.Cars[i].Position.X > 0 {
							models.Cars[i].Position.X -= 2  // Mueve hacia la izquierda
						} else {
							models.LaneMutex.Lock()
							models.LaneStatus[models.Cars[i].Lane] = false  // liberar el carril
							models.LaneMutex.Unlock()
							// Remover el carro de la lista
							models.Cars = append(models.Cars[:i], models.Cars[i+1:]...)
							exitIndex++  // Incrementa exitIndex cada vez que un carro sale.
						}
					} else {
						// Para carriles inferiores
						if models.Cars[i].Position.Y < 300 {
							models.Cars[i].Position.Y += 2  // Mueve hacia arriba
						} else if models.Cars[i].Position.X > 0 {
							models.Cars[i].Position.X -= 2  // Mueve hacia la izquierda
						} else {
							models.LaneMutex.Lock()
							models.LaneStatus[models.Cars[i].Lane] = false  // liberar el carril
							models.LaneMutex.Unlock()
							// Remover el carro de la lista
							models.Cars = append(models.Cars[:i], models.Cars[i+1:]...)
							exitIndex++  // Incrementa exitIndex cada vez que un carro sale.
						}
					}
				}
		}
		models.CarsMutex.Unlock()
	
		time.Sleep(16 * time.Millisecond)
	}
}

