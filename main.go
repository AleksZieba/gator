package main

import (
	"database/sql"
	"fmt"
	"os"

	config "github.com/AleksZieba/gator/internal/config"
	"github.com/AleksZieba/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db      *database.Queries
	pconfig *config.Config
}

func main() {
	var s state

	cfg := config.Read()
	s.pconfig = &cfg

	db, err := sql.Open("postgres", s.pconfig.DbUrl)

	if err != nil {
		fmt.Printf("%e", err)
	}

	s.db = database.New(db)

	var c commands
	c.cmds = make(map[string]func(*state, command) error)

	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerGetUsers)
	c.register("agg", handleAgg)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("feeds", handlerFeeds)
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerFollowing))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("browse", middlewareLoggedIn(handleBrowse))

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	com := command{
		name: args[1],
		args: args[2:],
	}

	err = c.run(&s, com)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg = config.Read()
	// fmt.Println(cfg.CurrentUserName)
	os.Exit(0)
}
