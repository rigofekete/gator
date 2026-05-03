package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/rigofekete/gator/internal/app"
	"github.com/rigofekete/gator/internal/cmd"
	"github.com/rigofekete/gator/internal/config"
	"github.com/rigofekete/gator/internal/database"
	"github.com/rigofekete/gator/internal/handlers"
	"github.com/rigofekete/gator/internal/tui/model"
)

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

	newState := &app.State{
		DB:  dbQueries,
		Cfg: &cfg,
	}

	cmds := cmd.Commands{
		RegisteredCommands: map[string]func(*app.State, cmd.Command) error{},
	}

	cmds.Register("login", handlers.HandlerLogin)
	cmds.Register("register", handlers.HandlerRegister)
	cmds.Register("reset", handlers.HandlerReset)
	cmds.Register("users", handlers.HandlerListUsers)
	cmds.Register("agg", handlers.HandlerAgg)
	cmds.Register("addfeed", middlewareLoggedIn(handlers.HandlerAddFeed))
	cmds.Register("follow", middlewareLoggedIn(handlers.HandlerFollowFeed))
	cmds.Register("unfollow", middlewareLoggedIn(handlers.HandlerUnfollowFeed))
	cmds.Register("following", middlewareLoggedIn(handlers.HandlerFollowing))
	cmds.Register("browse", middlewareLoggedIn(handlers.HandlerBrowsePosts))
	cmds.Register("feeds", handlers.HandlerGetFeeds)

	if len(os.Args) >= 2 && os.Args[1] == "--tui" {
		model.RunTUI(&cfg)
		return
	}

	if len(os.Args) >= 3 && os.Args[1] == "--exec" {
		cmdName := os.Args[2]
		cmdArgs := os.Args[3:]
		err := cmds.Run(newState, cmd.Command{Name: cmdName, Args: cmdArgs})
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

	err = cmds.Run(newState, cmd.Command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

func middlewareLoggedIn(handler func(*app.State, cmd.Command, database.User) error) func(*app.State, cmd.Command) error {
	return func(s *app.State, cmd cmd.Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		return handler(s, cmd, user)
	}
}
