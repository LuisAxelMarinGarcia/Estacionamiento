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


func Lane(id int) {
	CarsMutex.Lock()
	Car := Car{
		ID:       id,
		Position: pixel.V(0, 300),  
		Lane:     -1,
		Parked:   false,
	}
	Cars = append(Cars, Car)
	CarsMutex.Unlock()

	
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


	Entrance <- true  
	rand.Seed(time.Now().UnixNano())
	lanes := rand.Perm(numLanes) 
	var lane int
	foundLane := false  
	LaneMutex.Lock()
	for _, l := range lanes {
		if !LaneStatus[l] {
			lane = l
			LaneStatus[l] = true  
			foundLane = true 
			break
		}
	}
	LaneMutex.Unlock()

	<-Entrance  
	if !foundLane { 
		CarsMutex.Lock()
		for i := range Cars {
			if Cars[i].ID == id {
				Cars[i].Position = pixel.V(0, 300)  
			}
		}
		CarsMutex.Unlock()
		return
	}

	
	CarsMutex.Lock()
	for i := range Cars {
		if Cars[i].ID == id {
			Cars[i].Lane = lane
		}
	}
	CarsMutex.Unlock()
}
