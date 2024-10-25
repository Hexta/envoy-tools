package diff

type LineMove struct {
	Line   string
	NewPos int
	OldPos int
}

type Changes struct {
	Added     []string
	Group     string
	Modified  map[string]string
	Removed   []string
	Reordered []*LineMove
}

func (c *Changes) Empty() bool {
	return len(c.Added) == 0 && len(c.Modified) == 0 && len(c.Removed) == 0
}
