package cmd

import (
	"os"

	"github.com/Hexta/envoy-tools/cmd/envoy-tools/cmd/cp"
	"github.com/Hexta/envoy-tools/internal/version"
	log "github.com/sirupsen/logrus"
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
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cp.Cmd)
	rootCmd.AddCommand(DocsCmd)
}
