package diff

import (
	"github.com/google/go-cmp/cmp"
)

type ClustersDiffOptions struct {
	IncludedClusters []string
}

func (c ClustersDiffOptions) GetIncludedClustersAsMap() map[string]struct{} {
	m := make(map[string]struct{})

	for _, cluster := range c.IncludedClusters {
		m[cluster] = struct{}{}
	}

	return m
}

func Clusters(a map[string]interface{}, b map[string]interface{}, opts ClustersDiffOptions) (*Changes, error) {
	added := make([]string, 0)
	removed := make([]string, 0)
	modified := make(map[string]string)

	includeClusters := opts.GetIncludedClustersAsMap()

	a = filterIncludedClusters(a, includeClusters)
	b = filterIncludedClusters(b, includeClusters)

	for clusterName := range a {
		if _, ok := b[clusterName]; !ok {
			removed = append(removed, clusterName)
		} else {
			diffStr := cmp.Diff(a[clusterName], b[clusterName])
			modified[clusterName] = diffStr
		}
	}

	for clusterName := range b {
		if _, ok := a[clusterName]; !ok {
			added = append(added, clusterName)
		}
	}

	return &Changes{Group: "clusters", Added: added, Removed: removed, Modified: modified}, nil
}

func filterIncludedClusters(clusters map[string]interface{}, includeClusters map[string]struct{}) map[string]interface{} {
	if len(includeClusters) == 0 {
		return clusters
	}

	filteredClusters := make(map[string]interface{}, len(clusters))

	for name := range includeClusters {
		if _, ok := clusters[name]; ok {
			filteredClusters[name] = clusters[name]
		}
	}

	return filteredClusters
}
