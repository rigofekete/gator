package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rigofekete/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
		return
	}
	fmt.Printf("Read config: %+v\n", cfg)

	newState := state{&cfg}


	cmds := commands{
		cmdsMap: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)


	userArgs := os.Args
	

	if len(userArgs) < 2 {
		log.Fatalf("\nat least a command name arg is needed")
	}

	cmdName := userArgs[1]
	cmdArgs := []string{}
	if len(userArgs) > 2 {
		cmdArgs = userArgs[2:]
	}
	
	cmdData := command{
		name: cmdName,
		args: cmdArgs,	
	}


	if f, exists := cmds.cmdsMap[cmdData.name]; exists == true {
		err = f(&newState, cmdData)
		if err != nil {
			log.Fatalf("error running command %s. Error: %v\n", cmdData.name, err)
		}
	} else {
		log.Fatalf("function name '%s' does not exist", cmdData.name)
	}

	
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return
	}

	fmt.Printf("Final config file: %+v\n", cfg)
}
