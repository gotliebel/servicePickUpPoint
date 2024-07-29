package infrastructure

import (
	"fmt"
	"homework-1/internal/constant"
	"homework-1/internal/service"
	"homework-1/pkg/hash"
	"log"
)

func processCommand(serv *service.Service, s *Synchronization) {
	defer s.Wg.Done()
	for {
		select {
		case sourceCommand := <-s.CommandChan:
			switch sourceCommand[len(sourceCommand)-1] {
			case constant.CommandHelp:
				fmt.Println(constant.HelpTxt)
			case constant.CommandAccept:
				execAccept(serv, sourceCommand[:len(sourceCommand)-1])
			case constant.CommandBack:
				execBack(serv, sourceCommand[:len(sourceCommand)-1])
			case constant.CommandPickUp:
				execPickup(serv, sourceCommand[:len(sourceCommand)-1])
			case constant.CommandList:
				execList(serv, sourceCommand[:len(sourceCommand)-1])
			case constant.CommandReturn:
				execReturn(serv, sourceCommand[:len(sourceCommand)-1])
			case constant.CommandTakebacks:
				execTakebacks(serv, sourceCommand[:len(sourceCommand)-1])
			case constant.CommandChangeNumRoutines:
				execChange(sourceCommand[:len(sourceCommand)-1], s)
			case constant.CustomExit:
				return
			default:
				log.Printf(constant.WrongCmdTextTemplate, sourceCommand[len(sourceCommand)-1])
			}
			hash.GenerateHash()
		}
	}
}
