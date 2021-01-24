package commands

import "log"

type ClientGenerateCmd struct {
}

func (s *ClientGenerateCmd) Execute(args []string) error {
	log.Println("execute client generate")
	return nil
}
