package models

import (
	"math/rand"
	"time"
	"sync"
	"github.com/faiface/pixel"
)

const (
	numLanes  = 20
)

const (
	LaneWidth = 150.0
)

var (
	LaneStatus [numLanes]bool 
	Cars       []Car
	CarsMutex  sync.Mutex
	LaneMutex  sync.Mutex
	Entrance   = make(chan bool, 1)
)


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

func WaitForPosition(id int, targetX float64) {
    for {
        carPos := FindCarPosition(id)
        if carPos.X >= targetX {
            break
        }
        time.Sleep(16 * time.Millisecond)
    }
}

func FindAvailableLane() (int, bool) {
    LaneMutex.Lock()
    defer LaneMutex.Unlock()
    rand.Seed(time.Now().UnixNano())
    lanes := rand.Perm(numLanes)
    for _, l := range lanes {
        if !LaneStatus[l] {
            LaneStatus[l] = true
            return l, true
        }
    }
    return -1, false
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

func AssignLaneToCar(id int, lane int) {
    CarsMutex.Lock()
    defer CarsMutex.Unlock()
    for i := range Cars {
        if Cars[i].ID == id {
            Cars[i].Lane = lane
        }
    }
}

func Lane(id int) {
    CreateCar(id)
    WaitForPosition(id, 100)
    Entrance <- true
    lane, foundLane := FindAvailableLane()
    <-Entrance
    if !foundLane {
        ResetCarPosition(id)
        return
    }
    AssignLaneToCar(id, lane)
}
