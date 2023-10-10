package command

import (
	"errors"
	"fmt"
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

func FileFrom(labels []string) string {
	if len(labels) == 0 {
		return config.StdLabel
	}

	exp := regexp.MustCompile(config.RegexpNotLabel)
	for i, l := range labels {
		labels[i] = exp.ReplaceAllString(strings.ToLower(l), "")
	}

	tree := strings.Join(labels, config.StdSeparator)
	return strings.ToLower(tree)
}

func FindLines(content []byte, pattern string) [][]byte {
	regex := regexp.MustCompile(fmt.Sprintf(config.RegexpLine, regexp.QuoteMeta(pattern)))
	return regex.FindAll(content, -1)
}

func DeleteLines(content []byte, pattern string) []byte {
	regexPattern := regexp.MustCompile(fmt.Sprintf(config.RegexpLine, regexp.QuoteMeta(pattern)))
	return regexPattern.ReplaceAll(content, []byte(""))
}
