package cds

import (
	"context"
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

var diffCmdOpts = struct {
	Clusters []string
	Indent   int
	Stats    bool
}{}

var diffCmd = &cobra.Command{
	Use:   "diff IP:PORT IP:PORT",
	Short: "Compare Envoy CDS configuration from two Envoy instances",
	Args:  cobra.ExactArgs(2),
	Example: `# Diff all clusters
$ envoy-tools cp cds diff 127.0.0.1:18000 127.0.0.1:18001

# Diff specific clusters
$ envoy-tools cp cds diff 127.0.0.1:18000 127.0.0.1:18001 -c cluster-1 -c cluster-2
`,
	Run: diffCmdRunFunc,
}

func diffCmdRunFunc(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	urls := args

	grpcCallOptions := []grpc.CallOption{grpc.MaxCallRecvMsgSize(config.CpCmdGlobalOptions.MaxGrpcMessageSize)}
	grpcDialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cmList := fetchAllClusters(ctx, urls, grpcCallOptions, grpcDialOptions)

	clsDiffOpts := diff.ClustersDiffOptions{
		IncludedClusters: diffCmdOpts.Clusters,
	}

	changes, err := diff.Clusters(cmList[0], cmList[1], clsDiffOpts)
	if err != nil {
		log.WithError(err).Fatal("Failed to diff clusters")
	}

	diffStr := diff.FormatChanges(changes, diff.FormatOptions{Indent: diffCmdOpts.Indent, StatsOnly: diffCmdOpts.Stats})
	fmt.Println(diffStr)
}

func fetchAllClusters(ctx context.Context, urls []string, grpcCallOptions []grpc.CallOption, grpcDialOptions []grpc.DialOption) []map[string]interface{} {
	cmList := make([]map[string]interface{}, 2)

	var wg sync.WaitGroup
	for idx, url := range urls {
		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()
			xdsClient := util.NewXDSClient(url, grpcCallOptions, grpcDialOptions, config.CpCmdGlobalOptions.NodeID)
			cm, err := util.FetchClustersAsMap(ctx, xdsClient)
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
	diffCmd.Flags().IntVarP(&diffCmdOpts.Indent, "indent", "i", 4, "Indentation level")
	diffCmd.Flags().StringSliceVarP(&diffCmdOpts.Clusters, "cluster", "c", []string{}, "Cluster name")
	diffCmd.Flags().BoolVarP(&diffCmdOpts.Stats, "stats", "s", false, "Display stats only")
}
