package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/AleksZieba/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login expects a username")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	if err := s.pconfig.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Println("user set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("register expects a username")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})

	if err != nil {
		return err
	}

	if err = s.pconfig.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("user created and loged in: %v", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("this command takes no arguments")
	}
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return err
	}
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("this command takes no arguments")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.pconfig.CurrentUserName {
			fmt.Printf("%v (current)\n", user.Name)
			continue
		}
		fmt.Println(user.Name)
	}
	return nil
}

func handleAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("command should take a duration as an argument")
	}
	fmt.Printf("Collecting feeds every %v\n", cmd.args[0])

	timeBtwReq, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBtwReq)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	updatedfeed, err := s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     time.Now(),
		ID:            nextFeed.ID,
	})
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), updatedfeed.Url)
	if err != nil {
		return err
	}

	fmt.Println(feed.Channel.Title)
	for _, item := range feed.Channel.Items {
		pubTime := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			pubTime = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: pubTime,
			FeedID:      updatedfeed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Println(err)
			continue
		}
	}
	fmt.Printf("Feed %s collected, %v posts found\n", feed.Channel.Title, len(feed.Channel.Items))
	return nil
}
