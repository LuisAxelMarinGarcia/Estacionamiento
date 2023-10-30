package models

import (
	"time"
)

const (
	LaneWidth = 150.0
)

var AllParked bool  // false si alguno no está estacionado, true si todos están estacionados
var ExitIndex int


func MoveCarsLogic() {

	for i := len(Cars) - 1; i >= 0; i-- {

		if Cars[i].Position.X < 100 && Cars[i].Lane == -1 {
			Cars[i].Position.X += 2  // Mueve más lentamente hacia la derecha
			if Cars[i].Position.X > 100 {  // Asegura que el carro se detenga en X=100
				Cars[i].Position.X = 100
			}
		} else if Cars[i].Position.X == 100 && Cars[i].Lane == -1 {
			// La lógica de entrada existente...
		} else if Cars[i].Lane != -1 && !Cars[i].Parked {
			var targetX, targetY float64
			if Cars[i].Lane < 4 {
				// Carriles superiores
				targetX = 100.0 + float64(Cars[i].Lane)*LaneWidth + LaneWidth/2
				targetY = 350 + (500-350)/2  // Centro vertical de los carriles superiores
			} else {
				// Carriles inferiores
				targetX = 100.0 + float64(Cars[i].Lane-4)*LaneWidth + LaneWidth/2
				targetY = 100 + (250-100)/2  // Centro vertical de los carriles inferiores
			}
			if Cars[i].Position.X < targetX - 2 {
				Cars[i].Position.X += 2  // Mueve hacia la derecha
			} else if Cars[i].Position.Y < targetY - 2 && (targetX - Cars[i].Position.X) <= 2 {
				Cars[i].Position.Y += 2  // Mueve hacia arriba
			} else if Cars[i].Position.Y > targetY + 2 && (targetX - Cars[i].Position.X) <= 2 {
				Cars[i].Position.Y -= 2  // Mueve hacia abajo
			} else if (targetX - Cars[i].Position.X) <= 2 && (targetY - Cars[i].Position.Y) <= 2 {
				Cars[i].Position.X = targetX  // Ajusta la posición X a la posición objetivo
				Cars[i].Position.Y = targetY  // Ajusta la posición Y a la posición objetivo
				Cars[i].Parked = true  // Marca el carro como estacionado
				SetExitTime(&Cars[i])  // Establece el tiempo de salida del carro

				AllParked = true  // asume que todos están estacionados
				for _, car := range Cars {
					if !car.Parked {
						AllParked = false  // Si encuentra un carro que no está estacionado, establece allParked a false
						break
					}
				}
			}
			} else if Cars[i].Parked && time.Now().After(Cars[i].ExitTime) && AllParked {
				if ExitIndex >= 0 && ExitIndex < len(Cars) && i == ExitIndex {
					if Cars[i].Lane < 4 {
					// Para carriles superiores
					if Cars[i].Position.Y > 300 {
						Cars[i].Position.Y -= 2  // Mueve hacia abajo
					} else if Cars[i].Position.X > 0 {
						Cars[i].Position.X -= 2  // Mueve hacia la izquierda
					} else {
						LaneMutex.Lock()
						LaneStatus[Cars[i].Lane] = false  // liberar el carril
						LaneMutex.Unlock()
						// Remover el carro de la lista
						Cars = append(Cars[:i], Cars[i+1:]...)
						ExitIndex++  // Incrementa exitIndex cada vez que un carro sale.
					}
				} else {
					// Para carriles inferiores
					if Cars[i].Position.Y < 300 {
						Cars[i].Position.Y += 2  // Mueve hacia arriba
					} else if Cars[i].Position.X > 0 {
						Cars[i].Position.X -= 2  // Mueve hacia la izquierda
						} else {
							LaneMutex.Lock()
							LaneStatus[Cars[i].Lane] = false  // liberar el carril
							LaneMutex.Unlock()
							// Remover el carro de la lista
							Cars = append(Cars[:i], Cars[i+1:]...)
							ExitIndex = 0  // Restablece ExitIndex a 0 cada vez que un carro sale
						}
				}

	}
}

}

}