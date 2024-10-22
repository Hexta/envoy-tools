package diff

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/qri-io/deepdiff"
)

func Clusters(ctx context.Context, a map[string]interface{}, b map[string]interface{}) (*Changes, error) {
	dd := deepdiff.New()

	added := make([]string, 0)
	removed := make([]string, 0)
	modified := make(map[string]string)

	for clusterName := range a {
		if _, ok := b[clusterName]; !ok {
			removed = append(removed, clusterName)
		} else {
			diffs, stats, err := dd.StatDiff(ctx, a[clusterName], b[clusterName])
			if err != nil {
				return nil, err
			}

			if stats.NodeChange() == 0 || (diffs.Len() == 1 && diffs[0].Deltas == nil) {
				continue
			}

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
