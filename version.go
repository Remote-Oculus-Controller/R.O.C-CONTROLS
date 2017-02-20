package roc

import "fmt"

type version struct {
	Major, Minor, Patch int
	Label               string
	Name                string
}

//Version build
var Version = version{2, 0, 2, "dev", "Beta ssl v1"}

//Build definition
var Build string

func (v version) String() string {
	if v.Label != "" {
		return fmt.Sprintf("Roll version %d.%d.%d-%s \"%s\"\nGit commit hash: %s", v.Major, v.Minor, v.Patch, v.Label, v.Name, Build)
	}
	return fmt.Sprintf("Roll version %d.%d.%d \"%s\"\nGit commit hash: %s", v.Major, v.Minor, v.Patch, v.Name, Build)
}
