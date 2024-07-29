package infrastructure

import (
	"fmt"
	"homework-1/internal/constant"
	"homework-1/internal/service"
	"log"
)

func ControlRoutines(serv *service.Service, s *Synchronization) {
	defer s.Wg.Done()
	for {
		select {
		case val := <-s.NumRoutinesChan:
			if val < 1 {
				log.Printf("num routines is less than 1")
				continue
			}
			if val < s.NumRoutines {
				for i := 0; i < s.NumRoutines-val; i++ {
					s.CommandChan <- []string{constant.CustomExit}
				}
				fmt.Printf("changed number routines from %d to %d\n", s.NumRoutines, val)
				s.NumRoutines = val
				continue
			}
			if val > s.NumRoutines {
				for i := 0; i < val-s.NumRoutines; i++ {
					s.Wg.Add(1)
					go processCommand(serv, s)
				}
				fmt.Printf("changed number routines from %d to %d\n", s.NumRoutines, val)
				s.NumRoutines = val
				continue
			}
			fmt.Printf("number of routines is the same: %d\n", s.NumRoutines)
		case <-s.ExitChan:
			return
		}
	}
}
