package models

import (
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
)

const (
	numLanes  = 8
)


var (
	LaneStatus [numLanes]bool // false si está libre, true si está ocupado
	Cars       []Car
	CarsMutex  sync.Mutex
	LaneMutex  sync.Mutex
	Entrance   = make(chan bool, 1)  // canal para controlar el acceso a la entrada, ahora inicializado y con buffer
)

type Car struct {
	ID           int
	Position     pixel.Vec
	PreviousPosition pixel.Vec  // Nuevo campo para rastrear la posición anterior
	Lane         int
	Parked       bool
	ExitTime     time.Time
}


func SetExitTime(car *Car) {
	rand.Seed(time.Now().UnixNano())
	exitIn := time.Duration(rand.Intn(10)+1) * time.Second  // por ejemplo, entre 1 y 10 segundos
	car.ExitTime = time.Now().Add(exitIn)
}

func GetCars() []Car {
	return Cars
}

func CheckCollision(pos1, pos2 pixel.Vec) bool {
	distance := pos1.Sub(pos2).Len()
	return distance < 20.0  // Ajusta este valor según el tamaño de tus carros
}

func checkAllCollisions() {
	CarsMutex.Lock()
	defer CarsMutex.Unlock()

	for i := 0; i < len(Cars); i++ {
		for j := i + 1; j < len(Cars); j++ {
			if CheckCollision(Cars[i].Position, Cars[j].Position) {
				// Lógica para manejar la colisión (e.g., mover los autos, registrar el evento, etc.)
			}
		}
	}
}


func Carr(id int) {
	CarsMutex.Lock()
	Car := Car{
		ID:       id,
		Position: pixel.V(0, 300),  // Posición inicial fuera de la ventana
		Lane:     -1,
		Parked:   false,
	}
	Cars = append(Cars, Car)
	CarsMutex.Unlock()

	// Espera a que el auto llegue a la entrada antes de asignar un carril
	for {
		var carPos pixel.Vec
		CarsMutex.Lock()
		for _, car := range Cars {
			if car.ID == id {
				carPos = car.Position
				break
			}
		}
		CarsMutex.Unlock()
		if carPos.X >= 100 {
			break
		}
		time.Sleep(16 * time.Millisecond)
	}

	// Espera hasta que la entrada esté libre
	Entrance <- true  // Bloquea la entrada

	// Busca un carril libre de manera aleatoria
	rand.Seed(time.Now().UnixNano())
	lanes := rand.Perm(numLanes)  // Genera una permutación aleatoria de carriles
	var lane int
	foundLane := false  // Variable para verificar si encontró un carril
	LaneMutex.Lock()
	for _, l := range lanes {
		if !LaneStatus[l] {
			lane = l
			LaneStatus[l] = true  // Ocupa el carril
			foundLane = true  // Marca que encontró un carril
			break
		}
	}
	LaneMutex.Unlock()

	<-Entrance  // Libera la entrada para otro carro

	// Si no encontró un carril, regresa a la posición inicial y sale
	if !foundLane {  // Modificado para verificar la variable foundLane
		CarsMutex.Lock()
		for i := range Cars {
			if Cars[i].ID == id {
				Cars[i].Position = pixel.V(0, 300)  // Posición inicial
			}
		}
		CarsMutex.Unlock()
		return
	}

	// Actualiza el carril del carro pero no la posición
	CarsMutex.Lock()
	for i := range Cars {
		if Cars[i].ID == id {
			Cars[i].Lane = lane
		}
	}
	CarsMutex.Unlock()
}
