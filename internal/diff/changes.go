package diff

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

type LineMove struct {
	Line   string
	OldPos int
	NewPos int
}

type Changes struct {
	Modified  map[string]string
	Added     []string
	Removed   []string
	Group     string
	Reordered []*LineMove
}

func (c *Changes) Empty() bool {
	return len(c.Added) == 0 && len(c.Modified) == 0 && len(c.Removed) == 0
}

type FormatOptions struct {
	Indent    int
	StatsOnly bool
}

func FormatChanges(changes *Changes, opts FormatOptions) string {
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

func formatReordered(sb *strings.Builder, indent int, changes *Changes) {
	_, _ = fmt.Fprintf(sb, "%sreordered\n", strings.Repeat(" ", indent))
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

func formatModified(sb *strings.Builder, indent int, changes *Changes, opts FormatOptions) {
	_, _ = fmt.Fprintf(sb, "%smodified\n", strings.Repeat(" ", indent))
	names := maps.Keys(changes.Modified)
	slices.Sort(names)

	for _, name := range names {
		_, err := fmt.Fprintf(sb, "%s%s\n", strings.Repeat(" ", 2*indent), name)
		if err != nil {
			continue
		}

		if opts.StatsOnly {
			continue
		}

		_, err = fmt.Fprintf(sb, "%s%s\n", strings.Repeat(" ", 3*indent), changes.Modified[name])
		if err != nil {
			continue
		}
	}
}

func formatRemoved(sb *strings.Builder, indent int, changes *Changes) {
	_, _ = fmt.Fprintf(sb, "%sremoved\n", strings.Repeat(" ", indent))
	slices.Sort(changes.Removed)

	for _, it := range changes.Removed {
		_, err := fmt.Fprintf(sb, "%s%s\n", strings.Repeat(" ", 2*indent), it)
		if err != nil {
			continue
		}
	}
}

func formatAdded(sb *strings.Builder, indent int, changes *Changes) {
	slices.Sort(changes.Added)

	_, _ = fmt.Fprintf(sb, "%sadded\n", strings.Repeat(" ", indent))
	for _, it := range changes.Added {
		_, err := fmt.Fprintf(sb, "%s%s\n", strings.Repeat(" ", 2*indent), it)
		if err != nil {
			continue
		}
	}
}
