package config

import "github.com/Hexta/envoy-tools/internal/format"

var CommonOptions = struct {
	Format format.Format
}{
	Format: format.Format{
		Name: format.YAMLFormat,
	},
}
