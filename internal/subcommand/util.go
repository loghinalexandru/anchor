package subcommand

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	defaultLabel   = "root"
	defaultDir     = ".anchor"
	regexpNotLabel = `[^a-z0-9-\.]`
	regexpLine     = `(?im)^.+%s.+$`
)

var (
	ErrInvalidLabel = errors.New("invalid label passed")
)

func validate(labels []string) error {
	exp := regexp.MustCompile(`^[a-z0-9-\.]+$`)

	for _, l := range labels {
		if !exp.MatchString(l) {
			return fmt.Errorf("%s: %w", l, ErrInvalidLabel)
		}
	}

	return nil
}

func rootDir() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(home, defaultDir), nil
}

func formatLabels(labels []string) string {
	if len(labels) == 0 {
		return defaultLabel
	}

	tree := strings.Join(labels, ".")
	tree = strings.ToLower(tree)

	exp := regexp.MustCompile(regexpNotLabel)
	return exp.ReplaceAllString(tree, "")
}

func confirmation(s string, in io.Reader) bool {
	reader := bufio.NewReader(in)
	retryMax := 3

	for retryMax > 0 {
		fmt.Printf("%s [y/n]: ", s)

		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			fmt.Println("Aborting...")
			return false
		}

		retryMax--
	}

	fmt.Println("Exceeded retry count")
	return false
}

func findLines(content []byte, pattern string) [][]byte {
	regex := regexp.MustCompile(fmt.Sprintf(regexpLine, regexp.QuoteMeta(pattern)))
	return regex.FindAll(content, -1)
}
