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
	Clusters []string
	Indent   int
}{}

var diffCmd = &cobra.Command{
	Use:   "diff IP:PORT IP:PORT",
	Short: "Compare Envoy CDS configuration from two Envoy instances",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		urls := args

		grpcCallOptions := []grpc.CallOption{grpc.MaxCallRecvMsgSize(config.CpCmdGlobalOptions.MaxGrpcMessageSize)}
		grpcDialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		cmList := fetchAllClusters(urls, grpcCallOptions, grpcDialOptions)

		clsDiffOpts := diff.ClustersDiffOptions{
			IncludedClusters: diffOptions.Clusters,
		}

		changes, err := diff.Clusters(cmList[0], cmList[1], clsDiffOpts)
		if err != nil {
			log.WithError(err).Fatal("Failed to diff clusters")
		}

		diff := diff.FormatChanges(changes, diffOptions.Indent)
		fmt.Println(diff)
	},
}

func fetchAllClusters(urls []string, grpcCallOptions []grpc.CallOption, grpcDialOptions []grpc.DialOption) []map[string]interface{} {
	cmList := make([]map[string]interface{}, 2)

	var wg sync.WaitGroup
	for idx, url := range urls {
		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()
			xdsClient := util.NewXDSClient(url, grpcCallOptions, grpcDialOptions, config.CpCmdGlobalOptions.NodeID)
			cm, err := util.FetchClustersAsMap(xdsClient)
			if err != nil {
				log.WithError(err).Fatal("Failed to fetch clusters")
			}
			cmList[idx] = cm
		}(idx, url)
	}
	wg.Wait()

	return cmList
}

func init() {
	diffCmd.Flags().IntVarP(&diffOptions.Indent, "indent", "i", 4, "Indentation level")
	diffCmd.Flags().StringSliceVarP(&diffOptions.Clusters, "cluster", "c", []string{}, "Clusters name")
}
