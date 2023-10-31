package models

import (
	"sync"
)


const (
	numLanes  = 20
	LaneWidth = 150.0
)


var (
	LaneStatus [numLanes]bool 
	Cars       []Car
	CarsMutex  sync.Mutex
	LaneMutex  sync.Mutex
	CarEnteringOrExiting bool  
)
