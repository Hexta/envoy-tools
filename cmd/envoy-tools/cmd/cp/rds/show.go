package rds

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
	Use:   "show IP:PORT [virtual host name]...",
	Short: "Show Envoy RDS configuration",
	Args:  cobra.MinimumNArgs(1),
	Example: `# Show all route configs
$ envoy-tools cp rds show 127.0.0.1:18000

# Show specific route configs
$ envoy-tools cp rds show 127.0.0.1:18000 route-config-1 route-config-2
`,
	Run: showCmdRunFunc,
}

func showCmdRunFunc(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	url := args[0]

	virtualHosts := make([]string, 0)
	if len(args) > 1 {
		virtualHosts = args[1:]
	}

	grpcCallOptions := []grpc.CallOption{grpc.MaxCallRecvMsgSize(config.CpCmdGlobalOptions.MaxGrpcMessageSize)}
	grpcDialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	xdsClient := util.NewXDSClient(url, grpcCallOptions, grpcDialOptions, config.CpCmdGlobalOptions.NodeID)
	cm, err := util.FetchRoutesAsMap(ctx, xdsClient)
	if err != nil {
		log.WithError(err).Fatal("Failed to fetch virtual hosts")
	}

	var data interface{} = cm
	if len(virtualHosts) > 0 {
		data = make(map[string]interface{})
		for _, name := range virtualHosts {
			if _, ok := cm[name]; !ok {
				log.Fatalf("Virtual host %s not found", name)
			}
			data.(map[string]interface{})[name] = cm[name]
		}
	}

	yaml, err := format.YAML(data)
	if err != nil {
		log.WithError(err).Fatal("Failed to format as YAML")
	}
	fmt.Println(yaml)
}
