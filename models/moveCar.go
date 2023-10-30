package models

import (
	"time"
)

var AllParked bool
var ExitIndex int

func MoveCarsLogic() {

	for i := len(Cars) - 1; i >= 0; i-- {

		if Cars[i].Position.X < 100 && Cars[i].Lane == -1 {
			Cars[i].Position.X += 2
			if Cars[i].Position.X > 100 {
				Cars[i].Position.X = 100
			}
		} else if Cars[i].Lane != -1 && !Cars[i].Parked {
			var targetX, targetY float64
			laneWidth := 600.0 / 10 // Asumiendo que el ancho del estacionamiento es 600 y hay 10 conjuntos de carriles
			if Cars[i].Lane < 10 {
				targetX = 100.0 + float64(Cars[i].Lane)*laneWidth + laneWidth/2
				targetY = 350 + (500-350)/2
			} else {
				targetX = 100.0 + float64(Cars[i].Lane-10)*laneWidth + laneWidth/2
				targetY = 100 + (250-100)/2
			}
			if Cars[i].Position.X < targetX - 2 {
				Cars[i].Position.X += 2  
			} else if Cars[i].Position.Y < targetY - 2 && (targetX - Cars[i].Position.X) <= 2 {
				Cars[i].Position.Y += 2  
			} else if Cars[i].Position.Y > targetY + 2 && (targetX - Cars[i].Position.X) <= 2 {
				Cars[i].Position.Y -= 2 
			} else if (targetX - Cars[i].Position.X) <= 2 && (targetY - Cars[i].Position.Y) <= 2 {
				Cars[i].Position.X = targetX  
				Cars[i].Position.Y = targetY  
				Cars[i].Parked = true  
				SetExitTime(&Cars[i])  

				AllParked = true  
				for _, car := range Cars {
					if !car.Parked {
						AllParked = false 
						break
					}
				}
			}
			}
			for i := len(Cars) - 1; i >= 0; i-- {
				if Cars[i].Parked && time.Now().After(Cars[i].ExitTime) {
					if ExitIndex >= 0 && ExitIndex < len(Cars) && i == ExitIndex {
					if Cars[i].Lane < 10 {  // Cambiado de 4 a 10 para abarcar todos los carriles inferiores
						if Cars[i].Position.Y > 300 {
							Cars[i].Position.Y -= 1
						} else if Cars[i].Position.X > 0 {
							Cars[i].Position.X -= 1
						} else {
							LaneMutex.Lock()
							LaneStatus[Cars[i].Lane] = false
							LaneMutex.Unlock()
	
							Cars = append(Cars[:i], Cars[i+1:]...)
							ExitIndex++
						}
					} else {
						if Cars[i].Position.Y < 300 {  // Cambiado de 300 a 550 para asegurar que los autos en los carriles superiores salgan completamente del carril antes de moverse a la izquierda.
							Cars[i].Position.Y += 1
						} else if Cars[i].Position.X > 0 {
							Cars[i].Position.X -= 1
						} else {
							LaneMutex.Lock()
							LaneStatus[Cars[i].Lane] = false
							LaneMutex.Unlock()
	
							Cars = append(Cars[:i], Cars[i+1:]...)
							ExitIndex = 0
						}
					}
					
				}
			}
			
		}
	}
}