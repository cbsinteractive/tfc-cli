package cmd

import "io"

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
	if _, err := c.w.Write([]byte(Version)); err != nil {
		return err
	}
	return nil
}
