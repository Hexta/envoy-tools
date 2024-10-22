package util

import (
	"context"
	"fmt"

	_ "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/health_check/event_sinks/file/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/retry/host/previous_hosts/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"
	cdspbv3 "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	rdspbv3 "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
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
	nodeId          string
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
		nodeId:          nodeId,
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

func (c *XDSClient) FetchClusters() (*discoveryv3.DiscoveryResponse, error) {
	return c.cdsClient.FetchClusters(
		context.Background(), &discoveryv3.DiscoveryRequest{Node: &corev3.Node{Id: c.nodeId}}, c.grpcCallOptions...,
	)
}

func (c *XDSClient) FetchRoutes() (*discoveryv3.DiscoveryResponse, error) {
	return c.rdsClient.FetchRoutes(
		context.Background(), &discoveryv3.DiscoveryRequest{Node: &corev3.Node{Id: c.nodeId}}, c.grpcCallOptions...,
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

			resourcesMap[mt.GetName()] = pbs.AsMap()
			if mt.GetName() == "" {
				fmt.Printf("OPA\n")
			}

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

			resourcesMap[mt.GetName()] = pbs.AsMap()
		}
	}

	return resourcesMap, nil
}

func FetchEnvoyClusters(client *XDSClient) (*discoveryv3.DiscoveryResponse, error) {
	err := client.Connect()

	if err != nil {
		return nil, fmt.Errorf("failed to connect to xDS API: %w", err)
	}

	return client.FetchClusters()
}
func FetchEnvoyRoutes(client *XDSClient) (*discoveryv3.DiscoveryResponse, error) {
	err := client.Connect()

	if err != nil {
		return nil, fmt.Errorf("failed to connect to xDS API: %w", err)
	}

	return client.FetchRoutes()
}
