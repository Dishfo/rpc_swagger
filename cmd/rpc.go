package main

import (
	"github.com/Dishfo/rpc_swagger/cmd/commands"
	flags "github.com/jessevdk/go-flags"
	"log"
)

func main() {
	parser := flags.NewParser(nil, flags.Default)

	genpar, err := parser.AddCommand("generate", "", "", &commands.GeneratorCommand{})
	if err != nil {
		log.Fatal("add command failed because of ", err.Error())
	}

	for _, command := range genpar.Commands() {
		log.Println("support cmd ", command.Name, command.ShortDescription, command.LongDescription)
	}

	_, err = parser.Parse()
	if err != nil {
		log.Fatal("wrong argument  ", err.Error())
	}

}
