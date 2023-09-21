package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	StdFileMode    = os.FileMode(0644)
	StdLabel       = "root"
	StdDir         = ".anchor"
	StdSeparator   = "."
	regexpNotLabel = `[^a-z0-9-]`
	regexpLabel    = `^[a-z0-9-]+$`
	regexpLine     = `(?im)^.+%s.+$`
)

var (
	ErrInvalidLabel = errors.New("invalid label passed")
)

func Validate(labels []string) error {
	exp := regexp.MustCompile(regexpLabel)
	for _, l := range labels {
		if !exp.MatchString(l) {
			return fmt.Errorf("%s: %w", l, ErrInvalidLabel)
		}
	}

	return nil
}

func RootDir() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(home, StdDir), nil
}

func FileFrom(labels []string) string {
	if len(labels) == 0 {
		return StdLabel
	}

	exp := regexp.MustCompile(regexpNotLabel)
	for i, l := range labels {
		labels[i] = exp.ReplaceAllString(l, "")
	}

	tree := strings.Join(labels, StdSeparator)
	return strings.ToLower(tree)
}

func FindLines(content []byte, pattern string) [][]byte {
	regex := regexp.MustCompile(fmt.Sprintf(regexpLine, regexp.QuoteMeta(pattern)))
	return regex.FindAll(content, -1)
}
