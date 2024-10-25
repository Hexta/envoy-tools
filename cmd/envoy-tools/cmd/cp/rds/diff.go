package rds

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/Hexta/envoy-tools/internal/config"
	"github.com/Hexta/envoy-tools/internal/diff"
	"github.com/Hexta/envoy-tools/internal/format"
	"github.com/Hexta/envoy-tools/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var diffCmdOpts = struct {
	Indent          int
	OutputFormat    string
	RouteConfigName string
	Stats           bool
	VirtualHosts    []string
}{}

var diffCmd = &cobra.Command{
	Use:   "diff IP:PORT IP:PORT",
	Short: "Compare Envoy CDS configuration from two Envoy instances",
	Args:  cobra.ExactArgs(2),
	Example: `# Diff all virtual hosts
$ envoy-tools cp rds diff 127.0.0.1:18000 127.0.0.1:18001

# Diff specific virtual hosts
$ envoy-tools cp rds diff 127.0.0.1:18000 127.0.0.1:18001 -r virtual-host-1 -r virtual-host-2
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

	cmList := fetchAllRoutes(ctx, urls, grpcCallOptions, grpcDialOptions)

	vhDiffOpts := diff.VirtualHostsDiffOptions{
		IncludedVirtualHosts: diffCmdOpts.VirtualHosts,
		RouteConfigName:      diffCmdOpts.RouteConfigName,
	}

	changes, err := diff.VirtualHosts(cmList[0], cmList[1], vhDiffOpts)
	if err != nil {
		log.WithError(err).Fatal("Failed to diff routes")
	}

	var diffStr string

	switch diffCmdOpts.OutputFormat {
	case "text":
		diffStr = format.ChangesAsText(changes, format.Options{Indent: diffCmdOpts.Indent, StatsOnly: diffCmdOpts.Stats})
	case "yaml":
		diffStr, err = format.YAML(changes)
		if err != nil {
			log.WithError(err).Fatal("Failed to format changes")
			os.Exit(1)
		}
	default:
		log.Fatalf("Unknown output format: %s", diffCmdOpts.OutputFormat)
	}

	fmt.Println(diffStr)
}

func fetchAllRoutes(ctx context.Context, urls []string, grpcCallOptions []grpc.CallOption, grpcDialOptions []grpc.DialOption) []map[string]interface{} {
	cmList := make([]map[string]interface{}, 2)

	var wg sync.WaitGroup
	for idx, url := range urls {
		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()
			xdsClient := util.NewXDSClient(url, grpcCallOptions, grpcDialOptions, config.CpCmdGlobalOptions.NodeID)
			cm, err := util.FetchRoutesAsMap(ctx, xdsClient)
			if err != nil {
				log.WithError(err).Fatal("Failed to fetch routes")
			}
			cmList[idx] = cm
		}(idx, url)
	}
	wg.Wait()

	return cmList
}

func init() {
	diffCmd.Flags().IntVarP(&diffCmdOpts.Indent, "indent", "i", 4, "Indentation level")
	diffCmd.Flags().StringSliceVarP(&diffCmdOpts.VirtualHosts, "virtualhost", "r", []string{}, "Virtual host name")
	diffCmd.Flags().BoolVarP(&diffCmdOpts.Stats, "stats", "s", false, "Display stats only")
	diffCmd.Flags().StringVar(&diffCmdOpts.RouteConfigName, "route-config-name", "default", "Route config name")
	diffCmd.Flags().StringVarP(&diffCmdOpts.OutputFormat, "output-format", "o", "text", "Output format (text, yaml)")
}
