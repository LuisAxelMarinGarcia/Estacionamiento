package models

import (
	"github.com/faiface/pixel"
	"time"
	"math/rand"
)

type Car struct {
	ID           int
	Position     pixel.Vec
	PreviousPosition pixel.Vec 
	Lane         int
	Parked       bool
	ExitTime     time.Time
	IsEntering bool 
    Teleporting        bool
    TeleportStartTime  time.Time
}

var CarChannel chan Car

func Init() {
	CarChannel = make(chan Car)
	go CarGenerator()
}

func CreateCar(id int) Car {
    CarsMutex.Lock()
    defer CarsMutex.Unlock()
    Car := Car{
        ID:       id,
        Position: pixel.V(0, 300),
        Lane:     -1,
        Parked:   false,
    }
    Cars = append(Cars, Car)
    return Car
}


func SetExitTime(car *Car) {
    rand.Seed(time.Now().UnixNano())
    exitIn := time.Duration(rand.Intn(5)+1) * time.Second  
    car.ExitTime = time.Now().Add(exitIn)
}


func GetCars() []Car {
	return Cars
}

func AssignLaneToCar(id int, lane int) {
    CarsMutex.Lock()
    defer CarsMutex.Unlock()
    for i := range Cars {
        if Cars[i].ID == id {
            Cars[i].Lane = lane
        }
    }
}

func ResetCarPosition(id int) {
    CarsMutex.Lock()
    defer CarsMutex.Unlock()
    for i := range Cars {
        if Cars[i].ID == id {
            Cars[i].Position = pixel.V(0, 300)
        }
    }
}

func FindCarPosition(id int) pixel.Vec {
    CarsMutex.Lock()
    defer CarsMutex.Unlock()
    for _, car := range Cars {
        if car.ID == id {
            return car.Position
        }
    }
    return pixel.Vec{}
}

func ParkCar(car *Car, targetX, targetY float64) {
	car.Position.X = targetX
	car.Position.Y = targetY
	car.Parked = true
	SetExitTime(car)
}


func removeCar(index int) {
	Cars = append(Cars[:index], Cars[index+1:]...)
}

func CarGenerator() {
	id := 0
	for {
		id++
		car := CreateCar(id)
		CarChannel <- car
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+500))  
	}
}
