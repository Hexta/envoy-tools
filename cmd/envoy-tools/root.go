package main

import (
	"fmt"
	"os"

	"github.com/Hexta/envoy-tools/pkg/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "envoy-tools",
	Short:   "Envoy toolbox",
	Version: version.Version(),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
