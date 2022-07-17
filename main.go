package main

import (
	"github.com/nonedotone/smtp-proxy/cmd"
	"github.com/nonedotone/smtp-proxy/version"
)

var (
	Version = "unknown"
	Commit  = "unknown"
)

func main() {
	version.New(Version, Commit)
	cmd.Execute()
}
