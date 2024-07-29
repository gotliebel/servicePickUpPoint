package infrastructure

import (
	"fmt"
	"homework-1/internal/constant"
	"os"
	"os/signal"
	"syscall"
)

func WaitForExitSignal(signals chan os.Signal) {
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	fmt.Println("got signal")
}

func WaitForExit(s *Synchronization) {
	for i := 0; i < s.NumRoutines; i++ {
		s.CommandChan <- []string{constant.CustomExit}
	}
	s.ExitChan <- struct{}{}
	s.Wg.Wait()
}
