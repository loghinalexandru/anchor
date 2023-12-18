package config

import (
	"os"
)

const (
	StdDir            = ".anchor"
	StdConfigFile     = "config/anchor.yaml"
	StdSyncMsg        = "Sync bookmarks"
	StdFileMode       = os.FileMode(0o666)
	StdLabel          = "root"
	StdLabelSeparator = "."
)
