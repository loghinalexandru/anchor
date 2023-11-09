package config

import (
	"os"
)

const (
	StdFileMode  = os.FileMode(0o666)
	StdLabel     = "root"
	StdSeparator = "."
)
