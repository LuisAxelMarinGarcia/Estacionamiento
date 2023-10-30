package models

import (
	"github.com/faiface/pixel"
	"time"
	"math/rand"
)

type Car struct {
	ID               int
	Position         pixel.Vec
	PreviousPosition pixel.Vec
	Lane             int
	Parked           bool
	ExitTime         time.Time
}
func SetExitTime(car *Car) {
	rand.Seed(time.Now().UnixNano())
	exitIn := time.Duration(rand.Intn(10)+1) * time.Second  
	car.ExitTime = time.Now().Add(exitIn)
}

func GetCars() []Car {
	return Cars
}