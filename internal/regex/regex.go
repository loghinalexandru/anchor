package regex

import (
	"fmt"
	"regexp"
)

const (
	regexpLine = "(?im)^.*%s.*$"
)

func MatchLines(content []byte, pattern string) [][]byte {
	regex := regexp.MustCompile(fmt.Sprintf(regexpLine, regexp.QuoteMeta(pattern)))
	return regex.FindAll(content, -1)
}

func MatchTitle(content []byte) string {
	titleMatch := regexp.MustCompile("<title>(?P<title>.*)</title>")
	match := titleMatch.FindSubmatch(content)

	if len(match) == 0 {
		return ""
	}

	return string(match[1])
}
