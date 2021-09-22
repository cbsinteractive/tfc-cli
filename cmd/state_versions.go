package cmd

import (
	"errors"
	"fmt"
	"io"
)

type StateVersionsCmd struct {
	r  Runner
	os dependencyCaller
	w  io.Writer
}

func NewStateVersionsCmd(os dependencyCaller, w io.Writer) *StateVersionsCmd {
	return &StateVersionsCmd{
		os: os,
		w:  w,
	}
}

func (c *StateVersionsCmd) Name() string {
	return "stateversions"
}

func (c *StateVersionsCmd) Init(args []string) error {
	if len(args) < 1 {
		return errors.New("no subcommand given")
	}
	runners := []Runner{
		NewStateVersionsCurrentCmd(c.os, c.w),
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

func (c *StateVersionsCmd) Run() error {
	return c.r.Run()
}
