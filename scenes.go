package main

import (
	"time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
)

var allParked bool  // false si alguno no está estacionado, true si todos están estacionados
var exitIndex int


func run() {
	entrance = make(chan bool, 1) // canal con buffer de 1
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
			laneMutex.Lock()
			allLanesOccupied := true
			for _, occupied := range laneStatus {
				if !occupied {
					allLanesOccupied = false
					break
				}
			}
			laneMutex.Unlock()
			
			if allLanesOccupied {
				break  // Detiene el spawn de nuevos carros si todos los carriles están ocupados
			}
	
			go car(i)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000) + 1000))  // Espera más tiempo antes de lanzar el próximo carro
		}
	}()
	
	
	for !win.Closed() {
		win.Clear(colornames.White)
		drawParkingLot(win, GetCars())
		win.Update()
	
		carsMutex.Lock()
		for i := len(cars) - 1; i >= 0; i-- {
			for j := range cars {
				if i != j {  // Asegúrate de no comparar el carro consigo mismo
					if checkCollision(cars[i].Position, cars[j].Position) {
						continue  // Si hay colisión, salta al siguiente ciclo del bucle
					}
				}
			}
	
			if cars[i].Position.X < 100 && cars[i].Lane == -1 {
				cars[i].Position.X += 2  // Mueve más lentamente hacia la derecha
				if cars[i].Position.X > 100 {  // Asegura que el carro se detenga en X=100
					cars[i].Position.X = 100
				}
			} else if cars[i].Position.X == 100 && cars[i].Lane == -1 {
				// La lógica de entrada existente...
			} else if cars[i].Lane != -1 && !cars[i].Parked {
				var targetX, targetY float64
				if cars[i].Lane < 4 {
					// Carriles superiores
					targetX = 100.0 + float64(cars[i].Lane)*laneWidth + laneWidth/2
					targetY = 350 + (500-350)/2  // Centro vertical de los carriles superiores
				} else {
					// Carriles inferiores
					targetX = 100.0 + float64(cars[i].Lane-4)*laneWidth + laneWidth/2
					targetY = 100 + (250-100)/2  // Centro vertical de los carriles inferiores
				}
				if cars[i].Position.X < targetX - 2 {
					cars[i].Position.X += 2  // Mueve hacia la derecha
				} else if cars[i].Position.Y < targetY - 2 && (targetX - cars[i].Position.X) <= 2 {
					cars[i].Position.Y += 2  // Mueve hacia arriba
				} else if cars[i].Position.Y > targetY + 2 && (targetX - cars[i].Position.X) <= 2 {
					cars[i].Position.Y -= 2  // Mueve hacia abajo
				} else if (targetX - cars[i].Position.X) <= 2 && (targetY - cars[i].Position.Y) <= 2 {
					cars[i].Position.X = targetX  // Ajusta la posición X a la posición objetivo
					cars[i].Position.Y = targetY  // Ajusta la posición Y a la posición objetivo
					cars[i].Parked = true  // Marca el carro como estacionado
					setExitTime(&cars[i])  // Establece el tiempo de salida del carro

					allParked = true  // asume que todos están estacionados
					for _, car := range cars {
						if !car.Parked {
							allParked = false  // Si encuentra un carro que no está estacionado, establece allParked a false
							break
						}
					}
				}
				} else if cars[i].Parked && time.Now().After(cars[i].ExitTime) && allParked && i == exitIndex {
					if cars[i].Lane < 4 {
						// Para carriles superiores
						if cars[i].Position.Y > 300 {
							cars[i].Position.Y -= 2  // Mueve hacia abajo
						} else if cars[i].Position.X > 0 {
							cars[i].Position.X -= 2  // Mueve hacia la izquierda
						} else {
							laneMutex.Lock()
							laneStatus[cars[i].Lane] = false  // liberar el carril
							laneMutex.Unlock()
							// Remover el carro de la lista
							cars = append(cars[:i], cars[i+1:]...)
							exitIndex++  // Incrementa exitIndex cada vez que un carro sale.
						}
					} else {
						// Para carriles inferiores
						if cars[i].Position.Y < 300 {
							cars[i].Position.Y += 2  // Mueve hacia arriba
						} else if cars[i].Position.X > 0 {
							cars[i].Position.X -= 2  // Mueve hacia la izquierda
						} else {
							laneMutex.Lock()
							laneStatus[cars[i].Lane] = false  // liberar el carril
							laneMutex.Unlock()
							// Remover el carro de la lista
							cars = append(cars[:i], cars[i+1:]...)
							exitIndex++  // Incrementa exitIndex cada vez que un carro sale.
						}
					}
				}
		}
		carsMutex.Unlock()
	
		time.Sleep(16 * time.Millisecond)
	}
}

