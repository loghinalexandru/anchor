package regex

import (
	"fmt"
	"regexp"
)

const (
	regexpLine      = `(?im)^.*%s.*$`
	regexpEndOfLine = `(?im)\s.%s.$`
)

func MatchLines(content []byte, pattern string) bool {
	regex := regexp.MustCompile(fmt.Sprintf(regexpEndOfLine, regexp.QuoteMeta(pattern)))
	return regex.Match(content)
}

func FindLines(content []byte, pattern string) [][]byte {
	regex := regexp.MustCompile(fmt.Sprintf(regexpLine, regexp.QuoteMeta(pattern)))
	return regex.FindAll(content, -1)
}

func FindTitle(content []byte) string {
	titleMatch := regexp.MustCompile("<title>(?P<title>.*)</title>")
	match := titleMatch.FindSubmatch(content)

	if len(match) == 0 {
		return ""
	}

	return string(match[1])
}
