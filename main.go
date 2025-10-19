package main

import (
	// "fmt"
	"log"
	"os"
	"database/sql"


	"github.com/rigofekete/gator/internal/config"
	"github.com/rigofekete/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	// fmt.Printf("Read config: %+v\n", cfg)

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error opening SQL database. Error: %v", err) 
	}
	defer db.Close()
	dbQueries := database.New(db)


	newState := &state{
		db: dbQueries,
		cfg: &cfg,
	}


	cmds := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	if len(os.Args) < 2 {
		log.Fatalf("\nat least a command name arg is needed")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:] 
	

	err = cmds.run(newState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

	
	// cfg, err = config.Read()
	// if err != nil {
	// 	log.Fatalf("Error reading config file: %v", err)
	// 	return
	// }
	//
	// fmt.Printf("Final config file: %+v\n", cfg)
}
