package envoy_tools

import (
	"fmt"
	"strings"

	"envoy-tools/pkg/util"

	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CmdOptions struct {
	Urls             *[]string
	DiffEntityString string
}

var cmdOptions = CmdOptions{}

var cpNextCmd = &cobra.Command{
	Use:   "cp",
	Short: "Compare Envoy Control Plane configs",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
var cpNextDiffCmd = &cobra.Command{
	Use:   "diff FILE FILE",
	Short: "Compare Envoy Control Plane configs",
	Run: func(cmd *cobra.Command, args []string) {
		cmdOptions.Urls = &args

		cmList := make([]map[string]interface{}, 0, 2)

		det, err := util.NewDiffEntityType(cmdOptions.DiffEntityString)
		if err != nil {
			log.Fatalf("Invalid diff entity: %v", err)
		}

		var resources *discoveryv3.DiscoveryResponse
		for _, url := range *cmdOptions.Urls {
			var err error
			switch *det {
			case util.CdsDiff:
				resources, err = util.FetchEnvoyClusters(url)
				if err != nil {
					log.Fatalf("Failed to fetch resources %v: %v", url, err)
				}
			case util.RdsDiff:
				resources, err = util.FetchEnvoyRoutes(url)
				if err != nil {
					log.Fatalf("Failed to fetch resources %v: %v", url, err)
				}
			}

			cm, err := util.DiscoveryResourcesAsMap(resources)

			if err != nil {
				log.Fatalf("Failed to convert clusters to map %v", err)
			}

			cmList = append(cmList, cm)
		}

		err = util.PrintDiffInterface(cmList[0], cmList[1], det)
		if err != nil {
			log.Fatalf("Failed to print diff: %v", err)
		}

	},
	Args: cobra.ExactArgs(2),
}

func init() {
	cpNextCmd.AddCommand(cpNextDiffCmd)

	cpNextDiffCmd.Flags().StringVarP(
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

	rootCmd.AddCommand(cpNextCmd)
}
