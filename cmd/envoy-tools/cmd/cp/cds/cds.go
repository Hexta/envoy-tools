package cds

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "cds",
	Short: "CDS tools",
}

func init() {
	Cmd.AddCommand(diffCmd)
	Cmd.AddCommand(showCmd)
}
