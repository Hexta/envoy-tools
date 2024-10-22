package diff

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/maps"
)

func VirtualHosts(a map[string]interface{}, b map[string]interface{}) (*Changes, error) {
	added := make([]string, 0)
	removed := make([]string, 0)
	modified := make(map[string]string)

	// TODO: check all route configs
	rtA := a["default"]
	rtB := b["default"]

	rtmA := rtA.(map[string]interface{})
	rtmB := rtB.(map[string]interface{})

	vhA, ok := rtmA["virtualHosts"]
	if !ok {
		return nil, fmt.Errorf("RDS config has no VirtualHosts field")
	}

	vhB, ok := rtmB["virtualHosts"]
	if !ok {
		return nil, fmt.Errorf("RDS config has no VirtualHosts field")
	}

	vhlA := vhA.([]interface{})
	vhlB := vhB.([]interface{})

	vhNamesA := make([]string, 0, len(vhlA))
	vhNamesB := make([]string, 0, len(vhlB))

	vhmA := make(map[string]interface{}, len(vhlA))
	vhmB := make(map[string]interface{}, len(vhlB))

	reorderedMap := make(map[string]*LineMove, 0)

	for _, vhInfoInterface := range vhlA {
		vhInfo := vhInfoInterface.(map[string]interface{})
		vhNameI, ok := vhInfo["name"]
		if !ok {
			return nil, fmt.Errorf("virtual host has no name")
		}
		vhName := vhNameI.(string)
		vhNamesA = append(vhNamesA, vhName)
		vhmA[vhName] = vhInfo
	}

	for _, vhInfoInterface := range vhlB {
		vhInfo := vhInfoInterface.(map[string]interface{})
		vhNameI, ok := vhInfo["name"]
		if !ok {
			return nil, fmt.Errorf("virtual host has no name")
		}
		vhName := vhNameI.(string)
		vhNamesB = append(vhNamesB, vhName)
		vhmB[vhName] = vhInfo
	}

	for idx, vhName := range vhNamesA {
		if _, ok := vhmB[vhName]; !ok {
			removed = append(removed, vhName)
			continue
		}

		diffStr := cmp.Diff(vhmA[vhName], vhmB[vhName])

		if diffStr == "" {
			continue
		}

		modified[vhName] = diffStr

		if idx < len(vhNamesB) && vhNamesB[idx] != vhName {
			reorderedMap[vhName] = &LineMove{Line: vhName, OldPos: idx}
		}
	}

	for idx, vhName := range vhNamesB {
		if _, ok := vhmA[vhName]; !ok {
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
