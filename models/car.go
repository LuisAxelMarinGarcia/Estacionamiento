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
    exitIn := time.Duration(rand.Intn(1)+1) * time.Second  
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
	if car.Position.X < targetX - 2 {
		car.Position.X += 2  
	} else if car.Position.Y < targetY - 2 && (targetX - car.Position.X) <= 2 {
		car.Position.Y += 2  
	} else if car.Position.Y > targetY + 2 && (targetX - car.Position.X) <= 2 {
		car.Position.Y -= 2 
	} else if (targetX - car.Position.X) <= 2 && (targetY - car.Position.Y) <= 2 {
		car.Position.X = targetX  
		car.Position.Y = targetY  
		car.Parked = true  
		SetExitTime(car)  
	}
}

func removeCar(index int) {
	Cars = append(Cars[:index], Cars[index+1:]...)
}