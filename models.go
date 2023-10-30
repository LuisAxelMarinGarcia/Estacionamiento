package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
)

const (
	numLanes  = 8
	laneWidth = 150.0
)


var (
	laneStatus [numLanes]bool // false si está libre, true si está ocupado
	cars       []Car
	carsMutex  sync.Mutex
	laneMutex  sync.Mutex
	entrance   = make(chan bool, 1)  // canal para controlar el acceso a la entrada, ahora inicializado y con buffer
)

type Car struct {
	ID           int
	Position     pixel.Vec
	PreviousPosition pixel.Vec  // Nuevo campo para rastrear la posición anterior
	Lane         int
	Parked       bool
	ExitTime     time.Time
}


func setExitTime(car *Car) {
	rand.Seed(time.Now().UnixNano())
	exitIn := time.Duration(rand.Intn(10)+1) * time.Second  // por ejemplo, entre 1 y 10 segundos
	car.ExitTime = time.Now().Add(exitIn)
}

func GetCars() []Car {
	return cars
}

func checkCollision(pos1, pos2 pixel.Vec) bool {
	distance := pos1.Sub(pos2).Len()
	return distance < 20.0  // Ajusta este valor según el tamaño de tus carros
}

func checkAllCollisions() {
	carsMutex.Lock()
	defer carsMutex.Unlock()

	for i := 0; i < len(cars); i++ {
		for j := i + 1; j < len(cars); j++ {
			if checkCollision(cars[i].Position, cars[j].Position) {
				// Lógica para manejar la colisión (e.g., mover los autos, registrar el evento, etc.)
			}
		}
	}
}


func car(id int) {
	carsMutex.Lock()
	car := Car{
		ID:       id,
		Position: pixel.V(0, 300),  // Posición inicial fuera de la ventana
		Lane:     -1,
		Parked:   false,
	}
	cars = append(cars, car)
	carsMutex.Unlock()

	// Espera a que el auto llegue a la entrada antes de asignar un carril
	for {
		var carPos pixel.Vec
		carsMutex.Lock()
		for _, car := range cars {
			if car.ID == id {
				carPos = car.Position
				break
			}
		}
		carsMutex.Unlock()
		if carPos.X >= 100 {
			break
		}
		time.Sleep(16 * time.Millisecond)
	}

	// Espera hasta que la entrada esté libre
	entrance <- true  // Bloquea la entrada

	// Busca un carril libre de manera aleatoria
	rand.Seed(time.Now().UnixNano())
	lanes := rand.Perm(numLanes)  // Genera una permutación aleatoria de carriles
	var lane int
	foundLane := false  // Variable para verificar si encontró un carril
	laneMutex.Lock()
	for _, l := range lanes {
		if !laneStatus[l] {
			lane = l
			laneStatus[l] = true  // Ocupa el carril
			foundLane = true  // Marca que encontró un carril
			break
		}
	}
	laneMutex.Unlock()

	<-entrance  // Libera la entrada para otro carro

	// Si no encontró un carril, regresa a la posición inicial y sale
	if !foundLane {  // Modificado para verificar la variable foundLane
		carsMutex.Lock()
		for i := range cars {
			if cars[i].ID == id {
				cars[i].Position = pixel.V(0, 300)  // Posición inicial
			}
		}
		carsMutex.Unlock()
		return
	}

	// Actualiza el carril del carro pero no la posición
	carsMutex.Lock()
	for i := range cars {
		if cars[i].ID == id {
			cars[i].Lane = lane
		}
	}
	carsMutex.Unlock()
}
