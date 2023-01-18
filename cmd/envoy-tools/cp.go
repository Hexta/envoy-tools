package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Hexta/envoy-tools/pkg/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CmdOptions struct {
	Urls               *[]string
	DiffEntityString   string
	NodeId             string
	MaxGrpcMessageSize int
}

var cmdOptions = CmdOptions{}

var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Compare Envoy Control Plane configs",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
var cpDiffCmd = &cobra.Command{
	Use:   "diff IP:PORT IP:PORT",
	Short: "Compare Envoy Control Plane configs",
	Run: func(cmd *cobra.Command, args []string) {
		cmdOptions.Urls = &args

		cmList := make([]map[string]interface{}, 2)

		det, err := util.NewDiffEntityType(cmdOptions.DiffEntityString)
		if err != nil {
			log.Fatalf("Invalid diff entity: %v", err)
		}

		var wg sync.WaitGroup

		for idx, url := range *cmdOptions.Urls {
			wg.Add(1)
			go func(idx int, url string) {
				defer wg.Done()
				cm, err := util.FetchDiscoveryResourcesAsMap(det, url, cmdOptions.MaxGrpcMessageSize, cmdOptions.NodeId)

				if err != nil {
					log.Fatalf("Failed to convert clusters to map %v", err)
				}

				cmList[idx] = cm
			}(idx, url)
		}
		wg.Wait()

		err = util.PrintDiffMap(cmList[0], cmList[1], det)
		if err != nil {
			log.Fatalf("Failed to print diff: %v", err)
		}

	},
	Args: cobra.ExactArgs(2),
}

func init() {
	cpCmd.AddCommand(cpDiffCmd)
	cpCmd.PersistentFlags().IntVar(&cmdOptions.MaxGrpcMessageSize, "max-grpc-message-size", 100*1024*1024, "Max size of gRPC message")

	cpCmd.PersistentFlags().StringVar(&cmdOptions.NodeId, "node-id", "", "Node id used in discovery requests")
	err := cpCmd.MarkPersistentFlagRequired("node-id")
	if err != nil {
		log.Fatalf("Failed to configure CLI: %v", err)
	}

	cpDiffCmd.Flags().StringVarP(
		&cmdOptions.DiffEntityString,
		"entity-type",
		"t",
		"all",
		fmt.Sprintf("(%s)", strings.Join([]string{
			util.CdsDiff.String(),
			util.RdsDiff.String(),
			util.AllDiff.String(),
		}, ",")),
	)

	rootCmd.AddCommand(cpCmd)
}
