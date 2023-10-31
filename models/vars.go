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
	Entrance   = make(chan bool, 1)
	AllParked bool
	ExitIndex int
	CarEnteringOrExiting bool  
	ControlMutex sync.Mutex  
)

