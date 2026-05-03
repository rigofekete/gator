package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/rigofekete/gator/internal/config"
	"github.com/rigofekete/gator/internal/database"
	"github.com/rigofekete/gator/internal/tui/model"
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

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error opening SQL database. Error: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	newState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("browse", middlewareLoggedIn(handlerBrowsePosts))
	cmds.register("feeds", handlerGetFeeds)

	if len(os.Args) >= 2 && os.Args[1] == "--tui" {
		model.RunTUI()
		return
	}

	if len(os.Args) >= 3 && os.Args[1] == "--exec" {
		cmdName := os.Args[2]
		cmdArgs := os.Args[3:]
		err := cmds.run(newState, command{Name: cmdName, Args: cmdArgs})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		return
	}

	if len(os.Args) < 2 {
		fmt.Printf("usage: gator <command> [args]\n       gator --tui\n\n")
		fmt.Println("commands:")
		fmt.Println("  register <name>")
		fmt.Println("  login <name>")
		fmt.Println("  reset")
		fmt.Println("  users")
		fmt.Println("  addfeed <name> <url>")
		fmt.Println("  follow <url>")
		fmt.Println("  unfollow <url>")
		fmt.Println("  following")
		fmt.Println("  browse [limit]")
		fmt.Println("  agg <duration>")
		fmt.Println("  feeds")
		fmt.Println("")
		fmt.Println("  --tui    launch the interactive TUI")
		fmt.Println("")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(newState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		return handler(s, cmd, user)
	}
}
