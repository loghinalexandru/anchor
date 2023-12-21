package label

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/loghinalexandru/anchor/internal/config"
)

var (
	ErrInvalidLabel = errors.New("invalid or missing label passed")
)

var notLabelRegexp = regexp.MustCompile(`([^a-z0-9-]|^$)`)

func Validate(labels []string) error {
	var result error
	for _, l := range labels {
		if notLabelRegexp.MatchString(l) {
			result = errors.Join(result, fmt.Errorf("%q: %w", l, ErrInvalidLabel))
		}
	}

	return result
}

func Filepath(labels []string) string {
	rootDir := config.RootDir()

	if len(labels) == 0 {
		return filepath.Join(rootDir, config.StdLabel)
	}

	for i, l := range labels {
		labels[i] = notLabelRegexp.ReplaceAllString(strings.ToLower(l), "")
	}

	filename := strings.Join(labels, config.StdLabelSeparator)
	return filepath.Join(rootDir, strings.ToLower(filename))
}
