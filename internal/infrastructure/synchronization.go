package infrastructure

import "sync"

const baseNumRoutines = 5

type Synchronization struct {
	Wg              sync.WaitGroup
	CommandChan     chan []string
	NumRoutinesChan chan int
	NumRoutines     int
	ExitChan        chan struct{}
}

func NewSync() *Synchronization {
	return &Synchronization{
		CommandChan:     make(chan []string),
		NumRoutinesChan: make(chan int),
		NumRoutines:     0,
		ExitChan:        make(chan struct{}),
	}
}
