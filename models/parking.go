package models

import (
	"time"
)

func MoveCarsLogic() {
	for i := len(Cars) - 1; i >= 0; i-- {
		if Cars[i].Position.X < 100 && Cars[i].Lane == -1 && !Cars[i].IsEntering {
			Cars[i].Position.X += 10  
			if Cars[i].Position.X > 100 {
				Cars[i].Position.X = 100
			}
		} else if Cars[i].Lane != -1 && !Cars[i].Parked {
			var targetX, targetY float64
			laneWidth := 600.0 / 10
			if Cars[i].Lane < 10 {
				targetX = 100.0 + float64(Cars[i].Lane)*laneWidth + laneWidth/2
				targetY = 400 + (500-350)/2
			} else {
				targetX = 100.0 + float64(Cars[i].Lane-10)*laneWidth + laneWidth/2
				targetY = 100 + (250-100)/2
			}
			ParkCar(&Cars[i], targetX, targetY)
		}
	}
	ExitCarLogic()
}




func CheckAllParked() bool {
	allParked := true  
	for _, car := range Cars {
		if !car.Parked {
			allParked = false  
			break
		}
	}
	return allParked
}

func ExitCarLogic() {
    for i := len(Cars) - 1; i >= 0; i-- {
        if Cars[i].Parked && time.Now().After(Cars[i].ExitTime) && !Cars[i].IsEntering {
            if !Cars[i].Teleporting {
                Cars[i].Teleporting = true
                Cars[i].TeleportStartTime = time.Now()
                Cars[i].Position.X = 50
                Cars[i].Position.Y = 400
            } else if time.Since(Cars[i].TeleportStartTime) >= time.Millisecond*500 {
                updateLaneStatus(Cars[i].Lane, false)
                removeCar(i)
            }
        }
    }
}


func updateLaneStatus(lane int, status bool) {
	LaneMutex.Lock()
	defer LaneMutex.Unlock()
	LaneStatus[lane] = status
}

