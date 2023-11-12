package config

import (
	"os"
)

const (
	StdDir        = ".anchor"
	StdConfigFile = "config/anchor.yaml"
	StdFileMode   = os.FileMode(0o666)
	StdLabel      = "root"
	StdSeparator  = "."
)
