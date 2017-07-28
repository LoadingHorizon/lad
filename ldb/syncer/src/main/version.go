package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var Version = "unknown"

func init() {
	kingpin.HelpFlag.Short('h')
	kingpin.Version(Version)
}
