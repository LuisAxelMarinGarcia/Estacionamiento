package models

import (
	"time"
)

var AllParked bool
var ExitIndex int


type Position struct {
	X, Y float64
}


func MoveCarsLogic() {
	for i := len(Cars) - 1; i >= 0; i-- {
		moveCarToInitialPosition(&Cars[i])
		moveCarToParkingPosition(&Cars[i])
		moveCarToExit(&Cars[i], i)
	}
}

func moveCarToInitialPosition(car *Car) {
	if car.Position.X < 100 && car.Lane == -1 {
		car.Position.X += 2
		if car.Position.X > 100 {
			car.Position.X = 100
		}
	}
}

func moveCarToParkingPosition(car *Car) {
	if car.Lane != -1 && !car.Parked {
		targetX, targetY := calculateTargetPosition(car)
		moveCar(car, targetX, targetY)
	}
}

func calculateTargetPosition(car *Car) (float64, float64) {
	laneWidth := 600.0 / 10
	var targetX, targetY float64
	if car.Lane < 10 {
		targetX = 100.0 + float64(car.Lane)*laneWidth + laneWidth/2
		targetY = 350 + (500-350)/2
	} else {
		targetX = 100.0 + float64(car.Lane-10)*laneWidth + laneWidth/2
		targetY = 100 + (250-100)/2
	}
	return targetX, targetY
}

func moveCar(car *Car, targetX, targetY float64) {
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
		checkAllCarsParked()
	}
}

func checkAllCarsParked() {
	AllParked = true
	for _, car := range Cars {
		if !car.Parked {
			AllParked = false
			break
		}
	}
}

func moveCarToExit(car *Car, index int) {
	if car.Parked && time.Now().After(car.ExitTime) && AllParked {
		if ExitIndex >= 0 && ExitIndex < len(Cars) && index == ExitIndex {
			handleCarExiting(car, index)
		}
	}
}

func handleCarExiting(car *Car, index int) {
	if car.Lane < 10 {
		exitLowerLane(car, index)
	} else {
		exitUpperLane(car, index)
	}
}

func exitLowerLane(car *Car, index int) {
	if car.Position.Y > 300 {
		car.Position.Y -= 2
	} else if car.Position.X > 0 {
		car.Position.X -= 2
	} else {
		updateLaneStatus(car.Lane, false)
		removeCar(index)
		ExitIndex++
	}
}

func exitUpperLane(car *Car, index int) {
	if car.Position.Y < 300 {
		car.Position.Y += 2
	} else if car.Position.X > 0 {
		car.Position.X -= 2
	} else {
		updateLaneStatus(car.Lane, false)
		removeCar(index)
		ExitIndex = 0
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