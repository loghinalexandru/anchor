package subcommand

import (
	"errors"
	"regexp"
	"strings"
)

const (
	defaultRoot = "root"
)

var (
	flattenRep      = strings.NewReplacer(",", "", " ", ".")
	ErrInvalidLabel = errors.New("invalid label passed")
)

func flatten(labels string) (string, error) {
	if labels == "" {
		return defaultRoot, nil
	}

	fileName := flattenRep.Replace(labels)
	if ok, _ := regexp.MatchString(`^[a-z0-9\.]+$`, fileName); !ok {
		return "", ErrInvalidLabel
	}

	return fileName, nil
}
