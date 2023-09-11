package subcommand

import (
	"errors"
	"regexp"
	"strings"
)

const (
	defaultTree = "root"
)

var (
	flattenRep      = strings.NewReplacer(",", "", " ", ".")
	ErrInvalidLabel = errors.New("invalid label passed")
)

func flattenWithValidation(labels string) (string, error) {
	if labels == "" {
		return defaultTree, nil
	}

	fileName := flattenRep.Replace(labels)
	if ok, _ := regexp.MatchString(`^[a-z0-9\.]+$`, fileName); !ok {
		return "", ErrInvalidLabel
	}

	return fileName, nil
}

func flatten(labels string) string {
	if labels == "" {
		return defaultTree
	}

	fileName := strings.TrimSpace(labels)
	fileName = strings.ToLower(fileName)
	fileName = flattenRep.Replace(fileName)

	return fileName
}
