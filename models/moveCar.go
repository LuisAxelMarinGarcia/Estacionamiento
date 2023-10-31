package models

import (
	"time"
	"sync"
)

var AllParked bool
var ExitIndex int
var CarEnteringOrExiting bool  
var ControlMutex sync.Mutex  

func MoveCarsLogic() {
	for i := len(Cars) - 1; i >= 0; i-- {
		if Cars[i].Position.X < 100 && Cars[i].Lane == -1 && !Cars[i].IsEntering {
			Cars[i].Position.X += 2
			if Cars[i].Position.X > 100 {
				Cars[i].Position.X = 100
			}
		} else if Cars[i].Lane != -1 && !Cars[i].Parked {
			if !Cars[i].IsEntering { 
				var targetX, targetY float64
				laneWidth := 600.0 / 10
				if Cars[i].Lane < 10 {
					targetX = 100.0 + float64(Cars[i].Lane)*laneWidth + laneWidth/2
					targetY = 350 + (500-350)/2
				} else {
					targetX = 100.0 + float64(Cars[i].Lane-10)*laneWidth + laneWidth/2
					targetY = 100 + (250-100)/2
				}
				ParkCar(&Cars[i], targetX, targetY)
			}
		}
	}

	ExitCarLogic()
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
			if ExitIndex >= 0 && ExitIndex < len(Cars) && i == ExitIndex {
				if Cars[i].Lane < 10 {
					if Cars[i].Position.Y > 300 {
						Cars[i].Position.Y -= 6 
					} else if Cars[i].Position.X > 0 {
						Cars[i].Position.X -= 6 
					} else {
						updateLaneStatus(Cars[i].Lane, false)
						removeCar(i)
						ExitIndex++
					}
				} else {
					if Cars[i].Position.Y < 300 {
						Cars[i].Position.Y += 6 
					} else if Cars[i].Position.X > 0 {
						Cars[i].Position.X -= 6 
					} else {
						updateLaneStatus(Cars[i].Lane, false)
						removeCar(i)
						ExitIndex = 0
					}
				}
			}
		}
	}
}




func updateLaneStatus(lane int, status bool) {
	LaneMutex.Lock()
	defer LaneMutex.Unlock()
	LaneStatus[lane] = status
}

func removeCar(index int) {
	Cars = append(Cars[:index], Cars[index+1:]...)
}