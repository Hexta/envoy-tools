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

var grpcCallOptions = []grpc.CallOption{
	grpc.MaxCallRecvMsgSize(100 * 1024 * 1024),
}

var grpcDialOptions = []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
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

func FetchEnvoyClusters(u string) (*discoveryv3.DiscoveryResponse, error) {
	conn, err := grpc.Dial(u, grpcDialOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to xDS API: %w", err)
	}
	defer conn.Close()

	client := cdspbv3.NewClusterDiscoveryServiceClient(conn)

	clusters, err := client.FetchClusters(
		context.Background(), &discoveryv3.DiscoveryRequest{Node: &corev3.Node{Id: "boltcp"}}, grpcCallOptions...,
	)

	return clusters, err
}
func FetchEnvoyRoutes(u string) (*discoveryv3.DiscoveryResponse, error) {
	conn, err := grpc.Dial(u, grpcDialOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to xDS API: %w", err)
	}
	defer conn.Close()

	client := rdspbv3.NewRouteDiscoveryServiceClient(conn)

	routes, err := client.FetchRoutes(
		context.Background(), &discoveryv3.DiscoveryRequest{Node: &corev3.Node{Id: "boltcp"}}, grpcCallOptions...,
	)

	return routes, err
}
