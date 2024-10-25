package cds

import (
	"fmt"

	"github.com/Hexta/envoy-tools/internal/config"
	"github.com/Hexta/envoy-tools/internal/format"
	"github.com/Hexta/envoy-tools/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var showCmd = &cobra.Command{
	Use:   "show IP:PORT [cluster name]...",
	Short: "Show Envoy CDS configuration",
	Args:  cobra.MinimumNArgs(1),
	Example: `# Show all clusters
$ envoy-tools cp cds show 127.0.0.1:18000

# Show specific clusters
$ envoy-tools cp cds show 127.0.0.1:18000 cluster1 cluster2
`,
	Run: showCmdRunFunc,
}

func showCmdRunFunc(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	url := args[0]

	clusters := make([]string, 0)
	if len(args) > 1 {
		clusters = args[1:]
	}

	grpcCallOptions := []grpc.CallOption{grpc.MaxCallRecvMsgSize(config.CpCmdGlobalOptions.MaxGrpcMessageSize)}
	grpcDialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	xdsClient := util.NewXDSClient(url, grpcCallOptions, grpcDialOptions, config.CpCmdGlobalOptions.NodeID)
	cm, err := util.FetchClustersAsMap(ctx, xdsClient)
	if err != nil {
		log.WithError(err).Fatal("Failed to fetch clusters")
	}

	var data interface{} = cm
	if len(clusters) > 0 {
		data = make(map[string]interface{})
		for _, clusterName := range clusters {
			if _, ok := cm[clusterName]; !ok {
				log.Fatalf("Cluster %s not found", clusterName)
			}
			data.(map[string]interface{})[clusterName] = cm[clusterName]
		}
	}

	yaml, err := format.YAML(data)
	if err != nil {
		log.WithError(err).Fatal("Failed to format as YAML")
	}
	fmt.Println(yaml)
}