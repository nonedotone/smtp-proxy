package version

import (
	"fmt"
)

// const (
// 	goVersionRegex = `go[0-9]{1,2}\.[0-9]{1,3}\.[0-9]{1,3}`
// )

var (
	version = "unknown"
	commit  = "unknown"
)

func New(_version, _commit string) {
	version = _version
	commit = _commit
}

func Version() string {
	if commit != "unknown" && len(commit) > 12 {
		commit = commit[:11]
	}
	//buildOS := runtime.GOOS + "/" + runtime.GOARCH
	return fmt.Sprintf("%s %s", version, commit)
}
