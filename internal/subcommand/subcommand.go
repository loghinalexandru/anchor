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

func formatWithValidation(labels string) (string, error) {
	if labels == "" {
		return defaultTree, nil
	}

	fileName := flattenRep.Replace(labels)
	if ok, _ := regexp.MatchString(`^[a-z0-9-\.]+$`, fileName); !ok {
		return "", ErrInvalidLabel
	}

	return fileName, nil
}

func format(labels string) string {
	if labels == "" {
		return defaultTree
	}

	exp := regexp.MustCompile(`[^a-z0-9-\.]`)
	return exp.ReplaceAllString(strings.ToLower(labels), "")
}
