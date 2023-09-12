package subcommand

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/peterbourgon/ff/v4"
)

const (
	regexpLine = `(?im)^.+%s.+$`
)

var (
	ErrInvalidLabel = errors.New("invalid label passed")
)

// Fix this random root as well and rethink this
func formatWithValidation(labels ff.Flag) (string, error) {
	if labels.GetValue() == "" {
		return "root", nil
	}

	rep := strings.NewReplacer(",", "", " ", ".")
	val := labels.GetValue()

	fileName := rep.Replace(val)
	if ok, _ := regexp.MatchString(`^[a-z0-9-\.]+$`, fileName); !ok {
		return "", ErrInvalidLabel
	}

	return fileName, nil
}

func findLines(content []byte, pattern string) [][]byte {
	regex := regexp.MustCompile(fmt.Sprintf(regexpLine, regexp.QuoteMeta(pattern)))
	return regex.FindAll(content, -1)
}
