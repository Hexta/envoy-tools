package diff

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/maps"
)

type VirtualHostsDiffOptions struct {
	RouteConfigName      string
	IncludedVirtualHosts []string
}

func (v VirtualHostsDiffOptions) GetIncludedVirtualHostsAsMap() map[string]struct{} {
	m := make(map[string]struct{})

	for _, virtualHost := range v.IncludedVirtualHosts {
		m[virtualHost] = struct{}{}
	}

	return m
}

func VirtualHosts(a map[string]interface{}, b map[string]interface{}, opts VirtualHostsDiffOptions) (*Changes, error) {
	added := make([]string, 0)
	removed := make([]string, 0)
	modified := make(map[string]string)

	// Extract route configs for both maps
	routeConfigA, err := extractRouteConfig(a, opts.RouteConfigName)
	if err != nil {
		return nil, err
	}

	routeConfigB, err := extractRouteConfig(b, opts.RouteConfigName)
	if err != nil {
		return nil, err
	}

	hostNamesA, hostsMapA, err := extractVirtualHosts(routeConfigA)
	if err != nil {
		return nil, err
	}

	hostNamesB, hostsMapB, err := extractVirtualHosts(routeConfigB)
	if err != nil {
		return nil, err
	}

	reorderedMap := make(map[string]*LineMove)
	includedHosts := opts.GetIncludedVirtualHostsAsMap()

	for idx, vhName := range hostNamesA {
		if !isVirtualHostIncluded(vhName, includedHosts) {
			continue
		}

		if _, ok := hostsMapB[vhName]; !ok {
			removed = append(removed, vhName)
			continue
		}

		diffStr := cmp.Diff(hostsMapA[vhName], hostsMapB[vhName])

		if diffStr == "" {
			continue
		}

		modified[vhName] = diffStr

		if idx < len(hostNamesB) && hostNamesB[idx] != vhName {
			reorderedMap[vhName] = &LineMove{Line: vhName, OldPos: idx}
		}
	}

	for idx, vhName := range hostNamesB {
		if !isVirtualHostIncluded(vhName, includedHosts) {
			continue
		}

		if _, ok := hostsMapA[vhName]; !ok {
			added = append(added, vhName)
			delete(reorderedMap, vhName)
			continue
		}

		if _, ok := reorderedMap[vhName]; ok {
			reorderedMap[vhName].NewPos = idx
		}
	}

	return &Changes{
		Group:     "virtual hosts",
		Added:     added,
		Removed:   removed,
		Modified:  modified,
		Reordered: maps.Values(reorderedMap),
	}, nil
}

// isVirtualHostIncluded checks if a virtual host is included based on filter options
func isVirtualHostIncluded(vhName string, includedHosts map[string]struct{}) bool {
	if len(includedHosts) == 0 {
		return true
	}

	_, ok := includedHosts[vhName]
	return ok
}

// extractRouteConfig retrieves a route config map from a given configuration map
func extractRouteConfig(configMap map[string]interface{}, routeConfigName string) (map[string]interface{}, error) {
	rt, ok := configMap[routeConfigName]
	if !ok {
		return nil, fmt.Errorf("route config %s not found", routeConfigName)
	}
	rtMap, ok := rt.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("route config %s is not a valid map, got %T", routeConfigName, rt)
	}
	return rtMap, nil
}

// extractVirtualHosts normalizes the virtual host list into a map and list of names
func extractVirtualHosts(routeConfig map[string]interface{}) ([]string, map[string]interface{}, error) {
	hostsI, ok := routeConfig["virtualHosts"]
	if !ok {
		return nil, nil, fmt.Errorf("RDS config has no VirtualHosts field")
	}

	hosts := hostsI.([]interface{})

	names := make([]string, 0, len(hosts))
	hostByName := make(map[string]interface{}, len(hosts))

	for _, hostI := range hosts {
		host := hostI.(map[string]interface{})
		nameI, ok := host["name"]
		if !ok {
			return nil, nil, fmt.Errorf("virtual host has no name")
		}
		name := nameI.(string)
		names = append(names, name)
		hostByName[name] = host
	}
	return names, hostByName, nil
}
