package util

import (
	"fmt"
	"strings"
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

func FormatChanges(changesList []*Changes, indent int) string {
	var sb strings.Builder

	for _, changes := range changesList {
		if changes.Empty() {
			continue
		}

		//TODO handle printf errors
		_, _ = fmt.Fprintf(&sb, "%s\n", changes.Group)

		if len(changes.Added) > 0 {
			formatAdded(&sb, indent, changes)
		}

		if len(changes.Removed) > 0 {
			formatRemoved(&sb, indent, changes)
		}

		if len(changes.Modified) > 0 {
			formatModified(&sb, indent, changes)
		}

		if len(changes.Reordered) > 0 {
			formatReordered(&sb, indent, changes)
		}
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
			//TODO handle printf errors
			continue
		}
	}
}

func formatModified(sb *strings.Builder, indent int, changes *Changes) {
	_, _ = fmt.Fprintf(sb, "%smodified\n", strings.Repeat(" ", indent))
	for name, diff := range changes.Modified {
		_, err := fmt.Fprintf(sb, "%s%s\n", strings.Repeat(" ", 2*indent), name)
		if err != nil {
			//TODO handle printf errors
			continue
		}

		_, err = fmt.Fprintf(sb, "%s%s\n", strings.Repeat(" ", 3*indent), diff)
		if err != nil {
			//TODO handle printf errors
			continue
		}
	}
}

func formatRemoved(sb *strings.Builder, indent int, changes *Changes) {
	_, _ = fmt.Fprintf(sb, "%sremoved\n", strings.Repeat(" ", indent))
	for _, it := range changes.Removed {
		_, err := fmt.Fprintf(sb, "%s%s\n", strings.Repeat(" ", 2*indent), it)
		if err != nil {
			//TODO handle printf errors
			continue
		}
	}
}

func formatAdded(sb *strings.Builder, indent int, changes *Changes) {
	_, _ = fmt.Fprintf(sb, "%sadded\n", strings.Repeat(" ", indent))
	for _, it := range changes.Added {
		_, err := fmt.Fprintf(sb, "%s%s\n", strings.Repeat(" ", 2*indent), it)
		if err != nil {
			//TODO handle printf errors
			continue
		}
	}
}
