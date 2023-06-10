package util

import (
	"errors"
	"fmt"
)

type DiffEntityType string

func NewDiffEntityType(name string) (*DiffEntityType, error) {
	var t DiffEntityType

	switch name {
	case CdsDiff.String():
		t = CdsDiff
	case RdsDiff.String():
		t = RdsDiff
	case AllDiff.String():
		t = AllDiff
	default:
		return nil, errors.New("unknown diff entity type")
	}

	return &t, nil
}

const (
	CdsDiff DiffEntityType = "cds"
	RdsDiff DiffEntityType = "rds"
	AllDiff DiffEntityType = "all"
)

func (t DiffEntityType) String() string {
	return string(t)
}

type EnvoyCPResource struct {
	Items   map[string]interface{}
	Version string
}

type EnvoyCPSnapshot struct {
	VersionMap interface{}
	Resources  []EnvoyCPResource
}

func PrintDiffMap(a map[string]interface{}, b map[string]interface{}, entityType *DiffEntityType) error {
	diff, err := FastDiffMap(a, b, entityType)

	if err != nil {
		return err
	}

	fmt.Print(diff)

	return nil
}

func FastDiffMap(a map[string]interface{}, b map[string]interface{}, entityType *DiffEntityType) (string, error) {
	diffs := make([]*Changes, 0)

	switch *entityType {
	case CdsDiff:
		clusterDiffs, err := DiffClusters(a, b)
		if err != nil {
			return "", err
		}
		diffs = append(diffs, clusterDiffs)

	case RdsDiff:
		vhDiffs, err := DiffVirtualHosts(a, b)
		if err != nil {
			return "", err
		}
		diffs = append(diffs, vhDiffs)
	}

	const indent = 4

	return FormatChanges(diffs, indent), nil
}
