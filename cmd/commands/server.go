package commands

import "log"

type ServerGenerateCmd struct {
}

func (s *ServerGenerateCmd) Execute(args []string) error {
	log.Println("execute server generate")
	return nil
}
