package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rigofekete/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
		return
	}
	fmt.Printf("Read config: %+v\n", cfg)

	newState := state{
		cfg: &cfg,
	}


	cmds := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)


	

	if len(os.Args) < 2 {
		log.Fatalf("\nat least a command name arg is needed")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:] 


	err = cmds.run(&newState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatalf("error running command %s. Error: %v\n", cmdName, err)
	}

	
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return
	}

	fmt.Printf("Final config file: %+v\n", cfg)
}
