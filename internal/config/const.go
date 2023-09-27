package config

import (
	"os"
)

const (
	StdFileMode    = os.FileMode(0644)
	StdLabel       = "root"
	StdDir         = ".anchor"
	StdSeparator   = "."
	RegexpNotLabel = `[^a-z0-9-]`
	RegexpLabel    = `^[a-z0-9-]+$`
	RegexpLine     = `(?im)^.+%s.+\n`
)
