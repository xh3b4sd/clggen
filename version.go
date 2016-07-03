package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c *clggen) InitVersionCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "version",
		Short: "Show current version of the binary.",
		Long:  "Show current version of the binary.",
		Run:   c.ExecVersionCmd,
	}

	return newCmd
}

func (c *clggen) ExecVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", c.Version)
}
