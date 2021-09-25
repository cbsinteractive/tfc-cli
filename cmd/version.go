package cmd

import (
	"fmt"
	"io"
)

type VersionCmd struct {
	w io.Writer
}

func NewVersionCmd(w io.Writer) *VersionCmd {
	return &VersionCmd{
		w: w,
	}
}

func (c *VersionCmd) Name() string {
	return "version"
}

func (c *VersionCmd) Init([]string) error {
	return nil
}

func (c *VersionCmd) Run() error {
	label := ""
	if ReleaseLabel != "" {
		label = fmt.Sprintf("-%s", ReleaseLabel)
	}
	if _, err := c.w.Write([]byte(fmt.Sprintf("%s.%s.%s%s", Major, Minor, Patch, label))); err != nil {
		return err
	}
	return nil
}
