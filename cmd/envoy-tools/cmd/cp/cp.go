package cp

import (
	"github.com/Hexta/envoy-tools/cmd/envoy-tools/cmd/cp/cds"
	"github.com/Hexta/envoy-tools/cmd/envoy-tools/cmd/cp/rds"
	"github.com/Hexta/envoy-tools/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "cp",
	Short: "Compare Envoy Control Plane configs",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	Cmd.AddCommand(cds.Cmd)
	Cmd.AddCommand(rds.Cmd)

	Cmd.PersistentFlags().IntVar(&config.CpCmdGlobalOptions.MaxGrpcMessageSize, "max-grpc-message-size", 100*1024*1024, "Max size of gRPC message")

	Cmd.PersistentFlags().StringVar(&config.CpCmdGlobalOptions.NodeID, "node-id", "", "Node id used in discovery requests")
	err := Cmd.MarkPersistentFlagRequired("node-id")
	if err != nil {
		log.Fatalf("Failed to configure CLI: %v", err)
	}
}
