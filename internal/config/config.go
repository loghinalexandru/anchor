package config

import (
	"os"
	"path"
	"time"
)

const (
	StdDirName        = ".anchor"
	StdConfigPath     = ".config/anchor/anchor.yaml"
	StdStorageKey     = "storage"
	StdSyncModeKey    = "sync"
	StdHttpTimeout    = 3 * time.Second
	StdLocationKey    = "path"
	StdSyncMsg        = "Sync bookmarks"
	StdFileMode       = os.FileMode(0o666)
	StdLabel          = "root"
	StdLabelSeparator = "."
)

func Filepath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Cannot open config path")
	}

	return path.Join(home, StdConfigPath)
}
