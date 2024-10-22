package rds

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "rds",
	Short: "RDS tools",
}

func init() {
	Cmd.AddCommand(diffCmd)
}
