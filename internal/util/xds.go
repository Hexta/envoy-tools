package util

import (
	"context"
	"fmt"

	"github.com/Hexta/envoy-tools/internal/config"
	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/health_check/event_sinks/file/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/retry/host/previous_hosts/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"
	cdspbv3 "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	rdspbv3 "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

var grpcDialOptions = []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
}

type XDSClient struct {
	cdsClient       cdspbv3.ClusterDiscoveryServiceClient
	rdsClient       rdspbv3.RouteDiscoveryServiceClient
	nodeID          string
	url             string
	conn            *grpc.ClientConn
	grpcCallOptions []grpc.CallOption
	grpcDialOptions []grpc.DialOption
}

func NewXDSClient(
	url string,
	callOptions []grpc.CallOption,
	dialOptions []grpc.DialOption,
	nodeId string,
) *XDSClient {
	return &XDSClient{
		grpcCallOptions: callOptions,
		grpcDialOptions: dialOptions,
		url:             url,
		nodeID:          nodeId,
	}
}

func NewXDSClientFromConfig(
	url string,
) *XDSClient {
	callOptions := []grpc.CallOption{grpc.MaxCallRecvMsgSize(config.CpCmdGlobalOptions.MaxGrpcMessageSize)}
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	nodeID := config.CpCmdGlobalOptions.NodeID

	return &XDSClient{
		grpcCallOptions: callOptions,
		grpcDialOptions: dialOptions,
		url:             url,
		nodeID:          nodeID,
	}
}

func (c *XDSClient) Connect() error {
	var err error
	c.conn, err = grpc.NewClient(c.url, grpcDialOptions...)
	if err != nil {
		return err
	}

	c.cdsClient = cdspbv3.NewClusterDiscoveryServiceClient(c.conn)
	c.rdsClient = rdspbv3.NewRouteDiscoveryServiceClient(c.conn)

	return nil
}

func (c *XDSClient) Close() error {
	return c.conn.Close()
}

func (c *XDSClient) FetchClusters(ctx context.Context) (*discoveryv3.DiscoveryResponse, error) {
	return c.cdsClient.FetchClusters(
		ctx, &discoveryv3.DiscoveryRequest{Node: &corev3.Node{Id: c.nodeID}}, c.grpcCallOptions...,
	)
}

func (c *XDSClient) FetchRoutes(ctx context.Context) (*discoveryv3.DiscoveryResponse, error) {
	return c.rdsClient.FetchRoutes(
		ctx, &discoveryv3.DiscoveryRequest{Node: &corev3.Node{Id: c.nodeID}}, c.grpcCallOptions...,
	)
}

func DiscoveryResourcesAsMap(clusters *discoveryv3.DiscoveryResponse) (map[string]interface{}, error) {
	resourcesMap := make(map[string]interface{})

	for _, resource := range clusters.GetResources() {
		m, err := resource.UnmarshalNew()
		if err != nil {
			return resourcesMap, err
		}

		switch mt := m.(type) {
		case *clusterv3.Cluster:
			buff, err := protojson.Marshal(resource)
			if err != nil {
				return resourcesMap, err
			}

			pbs := structpb.Struct{}
			err = protojson.Unmarshal(buff, &pbs)
			if err != nil {
				return resourcesMap, err
			}

			if mt.GetName() == "" {
				log.Errorf("Clusters name is empty")
			}
			resourcesMap[mt.GetName()] = pbs.AsMap()

		case *routev3.RouteConfiguration:
			buff, err := protojson.Marshal(resource)
			if err != nil {
				return resourcesMap, err
			}

			pbs := structpb.Struct{}
			err = protojson.Unmarshal(buff, &pbs)
			if err != nil {
				return resourcesMap, err
			}

			if mt.GetName() == "" {
				log.Errorf("Routes Configuration name is empty")
			}
			resourcesMap[mt.GetName()] = pbs.AsMap()
		}
	}

	return resourcesMap, nil
}

func FetchClusters(ctx context.Context, client *XDSClient) (*discoveryv3.DiscoveryResponse, error) {
	err := client.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to xDS API: %w", err)
	}

	return client.FetchClusters(ctx)
}

func FetchClustersAsMap(ctx context.Context, client *XDSClient) (map[string]interface{}, error) {
	clusters, err := FetchClusters(ctx, client)
	if err != nil {
		return nil, err
	}
	return DiscoveryResourcesAsMap(clusters)
}

func FetchRoutes(ctx context.Context, client *XDSClient) (*discoveryv3.DiscoveryResponse, error) {
	err := client.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to xDS API: %w", err)
	}

	return client.FetchRoutes(ctx)
}

func FetchRoutesAsMap(ctx context.Context, client *XDSClient) (map[string]interface{}, error) {
	routes, err := FetchRoutes(ctx, client)
	if err != nil {
		return nil, err
	}
	return DiscoveryResourcesAsMap(routes)
}
