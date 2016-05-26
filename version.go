package roc

import "fmt"

type version struct {
	Major, Minor, Patch int
	Label               string
	Name                string
}

var Version = version{0, 2, 1, "dev", "baby"}

var Build string

func (v version) String() string {
	if v.Label != "" {
		return fmt.Sprintf("Roll version %d.%d.%d-%s \"%s\"\nGit commit hash: %s", v.Major, v.Minor, v.Patch, v.Label, v.Name, Build)
	} else {
		return fmt.Sprintf("Roll version %d.%d.%d \"%s\"\nGit commit hash: %s", v.Major, v.Minor, v.Patch, v.Name, Build)
	}
}
