package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, exists := c.cmds[cmd.name]

	if !exists {
		return fmt.Errorf("Command not found")
	}
	if err := f(s, cmd); err != nil {
		return err
	}
	return nil
}
