package util

type EnvoyCPResource struct {
	Items   map[string]interface{}
	Version string
}

type EnvoyCPSnapshot struct {
	VersionMap interface{}
	Resources  []EnvoyCPResource
}
