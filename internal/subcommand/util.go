package subcommand

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	stdFileMode    = os.FileMode(0644)
	stdLabel       = "root"
	stdDir         = ".anchor"
	stdSeparator   = "."
	regexpNotLabel = `[^a-z0-9-]`
	regexpLabel    = `^[a-z0-9-]+$`
	regexpLine     = `(?im)^.+%s.+$`
)

var (
	ErrInvalidLabel = errors.New("invalid label passed")
)

func validate(labels []string) error {
	exp := regexp.MustCompile(regexpLabel)
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

	return filepath.Join(home, stdDir), nil
}

func fileFrom(labels []string) string {
	if len(labels) == 0 {
		return stdLabel
	}

	exp := regexp.MustCompile(regexpNotLabel)
	for i, l := range labels {
		labels[i] = exp.ReplaceAllString(l, "")
	}

	tree := strings.Join(labels, stdSeparator)
	return strings.ToLower(tree)
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

func lineCounter(r io.Reader) (int, error) {
	var res int
	buf := make([]byte, 32*1024)
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		res += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return res, nil
		case err != nil:
			return res, err
		}
	}
}
