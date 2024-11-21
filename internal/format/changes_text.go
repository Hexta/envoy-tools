package format

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Hexta/envoy-tools/internal/diff"
	"golang.org/x/exp/maps"
)

func ChangesAsText(changes *diff.Changes, opts Options) string {
	var sb strings.Builder

	if changes.Empty() {
		return ""
	}

	_, _ = fmt.Fprintf(&sb, "%s\n", changes.Group)

	if len(changes.Added) > 0 {
		formatAdded(&sb, opts.Indent, changes)
	}

	if len(changes.Removed) > 0 {
		formatRemoved(&sb, opts.Indent, changes)
	}

	if len(changes.Modified) > 0 {
		formatModified(&sb, opts.Indent, changes, opts)
	}

	if len(changes.Reordered) > 0 {
		formatReordered(&sb, opts.Indent, changes)
	}

	return sb.String()
}

func formatReordered(sb *strings.Builder, indent int, changes *diff.Changes) {
	writeLine(sb, "reordered", indent)

	for _, lineMove := range changes.Reordered {
		_, err := fmt.Fprintf(
			sb,
			"%s%s [%d] -> [%d]\n", strings.Repeat(" ", 2*indent), lineMove.Line, lineMove.OldPos, lineMove.NewPos,
		)
		if err != nil {
			continue
		}
	}
}

func formatModified(sb *strings.Builder, indent int, changes *diff.Changes, opts Options) {
	writeLine(sb, "modified", indent)

	names := maps.Keys(changes.Modified)
	slices.Sort(names)

	for _, name := range names {
		writeLine(sb, name, 2*indent)

		if opts.StatsOnly {
			continue
		}

		writeLine(sb, changes.Modified[name], 3*indent)
	}
}

func formatRemoved(sb *strings.Builder, indent int, changes *diff.Changes) {
	writeLine(sb, "removed", indent)

	slices.Sort(changes.Removed)
	writeLines(sb, changes.Removed, 2*indent)
}

func formatAdded(sb *strings.Builder, indent int, changes *diff.Changes) {
	writeLine(sb, "added", indent)

	slices.Sort(changes.Added)
	writeLines(sb, changes.Added, 2*indent)
}

// Helper to write a single indented line
func writeLine(sb *strings.Builder, content string, indent int) {
	_, _ = fmt.Fprintf(sb, "%s%s\n", strings.Repeat(" ", indent), content)
}

// Helper to write a list of lines with indentation
func writeLines(sb *strings.Builder, items []string, indent int) {
	for _, item := range items {
		writeLine(sb, item, indent)
	}
}
