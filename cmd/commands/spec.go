package commands

import "log"

type SpecGenerateCmd struct {
}

func (s *SpecGenerateCmd) Execute(args []string) error {
	log.Println("execute spec generate maybe")
	return nil
}
