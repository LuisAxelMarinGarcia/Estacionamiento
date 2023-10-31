package models

import (
	"math/rand"
	"time"
)

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
