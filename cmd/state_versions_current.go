package cmd

import (
	"errors"
	"fmt"
	"io"
)

type StateVersionsCurrentCmd struct {
	r  Runner
	os dependencyCaller
	w  io.Writer
}

func NewStateVersionsCurrentCmd(os dependencyCaller, w io.Writer) *StateVersionsCurrentCmd {
	return &StateVersionsCurrentCmd{
		os: os,
		w:  w,
	}
}

func (c *StateVersionsCurrentCmd) Name() string {
	return "current"
}

func (c *StateVersionsCurrentCmd) Init(args []string) error {
	if len(args) < 1 {
		return errors.New("no subcommand given")
	}
	runners := []Runner{
		NewStateVersionsCurrentGetOutputCmd(c.os, c.w),
	}
	subcommand := args[0]
	for _, r := range runners {
		if r.Name() == subcommand {
			if err := r.Init(args[1:]); err != nil {
				return err
			}
			c.r = r
			return nil
		}
	}
	return fmt.Errorf("unexpected subcommand: %s", subcommand)
}

func (c *StateVersionsCurrentCmd) Run() error {
	return c.r.Run()
}
