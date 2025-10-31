package main

import (
	"fmt"
	"log"
	"os"

	"github.com/allscorpion/gator/internal/commands"
	"github.com/allscorpion/gator/internal/config"
)

func main() {
	cnfg, err := config.Read();

	if err != nil {
		fmt.Print(err)
		return;
	}

	state := commands.State{
		Cnfg: &cnfg,
	}

	cmnds := commands.Commands{
		Options: map[string]func(*commands.State, commands.Command) error{},
	}

	cmnds.Register("login", commands.HandlerLogin)

	args := os.Args

	if len(args) < 2 {
		log.Fatal("a minimum of two arguments is required")
	}

	firstArg := args[1];
	restOfArgs := args[2:];

	cmd := commands.Command{
		Name: firstArg,
		Arguments: restOfArgs,
	}

	err = cmnds.Run(&state, cmd)

	if err != nil {
		log.Fatalf("%v", err)
	}
}