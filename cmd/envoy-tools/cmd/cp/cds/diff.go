package cds

import (
	"fmt"
	"sync"

	"github.com/Hexta/envoy-tools/internal/config"
	"github.com/Hexta/envoy-tools/internal/diff"
	"github.com/Hexta/envoy-tools/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var diffOptions = struct {
	Cluster string
	Indent  int
}{}

var diffCmd = &cobra.Command{
	Use:   "diff IP:PORT IP:PORT",
	Short: "Compare Envoy CDS configuration from two Envoy instances",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		urls := args

		cmList := make([]map[string]interface{}, 2)

		grpcCallOptions := []grpc.CallOption{grpc.MaxCallRecvMsgSize(config.CpCmdGlobalOptions.MaxGrpcMessageSize)}
		grpcDialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		var wg sync.WaitGroup
		for idx, url := range urls {
			wg.Add(1)
			go func(idx int, url string) {
				defer wg.Done()
				cm, err := fetchEnvoyClusters(url, grpcCallOptions, grpcDialOptions)
				if err != nil {
					log.WithError(err).Fatal("Failed to fetch clusters")
				}
				cmList[idx] = cm
			}(idx, url)
		}
		wg.Wait()

		changes, err := diff.Clusters(ctx, cmList[0], cmList[1])
		if err != nil {
			log.WithError(err).Fatal("Failed to diff clusters")
		}

		diff := diff.FormatChanges(changes, diffOptions.Indent)
		fmt.Println(diff)
	},
}

func fetchEnvoyClusters(url string, grpcCallOptions []grpc.CallOption, grpcDialOptions []grpc.DialOption) (map[string]interface{}, error) {
	xdsClient := util.NewXDSClient(url, grpcCallOptions, grpcDialOptions, config.CpCmdGlobalOptions.NodeID)

	resp, err := util.FetchEnvoyClusters(xdsClient)
	if err != nil {
		log.Fatalf("Failed to fetch clusters %v", err)
	}

	cm, err := util.DiscoveryResourcesAsMap(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert clusters to map %v", err)
	}

	return cm, nil
}

func init() {
	diffCmd.Flags().IntVarP(&diffOptions.Indent, "indent", "i", 4, "Indentation level")
	diffCmd.Flags().StringVarP(&diffOptions.Cluster, "cluster", "c", "", "Cluster name")
}
