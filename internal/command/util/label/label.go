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
	ErrInvalidLabel = errors.New("invalid label passed")
)

func Validate(labels []string) error {
	exp := regexp.MustCompile(config.RegexpLabel)
	for _, l := range labels {
		if !exp.MatchString(l) {
			return fmt.Errorf("%s: %w", l, ErrInvalidLabel)
		}
	}

	return nil
}

func Filepath(labels []string) string {
	rootDir := config.RootDir()

	if len(labels) == 0 {
		return filepath.Join(rootDir, config.StdLabel)
	}

	exp := regexp.MustCompile(config.RegexpNotLabel)
	for i, l := range labels {
		labels[i] = exp.ReplaceAllString(strings.ToLower(l), "")
	}

	filename := strings.Join(labels, config.StdSeparator)
	return filepath.Join(rootDir, strings.ToLower(filename))
}
