package regex

import (
	"fmt"
	"regexp"
)

const (
	regexpLine      = `(?im)^.*%s.*$`
	regexpEndOfLine = `(?im)\s.%s.$`
)

func MatchEndOfLines(content []byte, pattern string) bool {
	regex := regexp.MustCompile(fmt.Sprintf(regexpEndOfLine, regexp.QuoteMeta(pattern)))
	return regex.Match(content)
}

func FindLines(content []byte, pattern string) [][]byte {
	regex := regexp.MustCompile(fmt.Sprintf(regexpLine, regexp.QuoteMeta(pattern)))
	return regex.FindAll(content, -1)
}
