package models

import (
	"github.com/faiface/pixel"
	"time"
	"math/rand"
)

type Car struct {
	ID           int
	Position     pixel.Vec
	PreviousPosition pixel.Vec  // Nuevo campo para rastrear la posición anterior
	Lane         int
	Parked       bool
	ExitTime     time.Time
	IsEntering bool // Nuevo campo para rastrear si el carro está entrando
}


func SetExitTime(car *Car) {
    rand.Seed(time.Now().UnixNano())
    exitIn := time.Duration(rand.Intn(5)+1) * time.Second  // Ajustado para 1 a 5 segundos
    car.ExitTime = time.Now().Add(exitIn)
}


func GetCars() []Car {
	return Cars
}