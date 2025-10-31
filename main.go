package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/allscorpion/gator/internal/commands"
	"github.com/allscorpion/gator/internal/config"
	"github.com/allscorpion/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cnfg, err := config.Read();

	if err != nil {
		fmt.Print(err)
		return;
	}

	db, err := sql.Open("postgres", cnfg.Db_url)

	if err != nil {
		fmt.Print(err)
		return;
	}

	dbQueries := database.New(db)

	state := commands.State{
		Cnfg: &cnfg,
		Db: dbQueries,
	}

	cmnds := commands.Commands{
		Options: map[string]func(*commands.State, commands.Command) error{},
	}

	cmnds.Register("login", commands.HandlerLogin)
	cmnds.Register("register", commands.HandlerRegister)
	cmnds.Register("reset", commands.HandleReset)
	cmnds.Register("users", commands.HandleGetUsers)

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