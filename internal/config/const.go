package config

import (
	"os"
)

const (
	StdFileMode    = os.FileMode(0644)
	StdLabel       = "root"
	StdSeparator   = "."
	RegexpNotLabel = `[^a-z0-9-]`
	RegexpLabel    = `^[a-z0-9-]+$`
	RegexpLine     = `(?i).+%s.+ .+\n`
	RegexpURL      = `(?im)\s.%s.$`
)
