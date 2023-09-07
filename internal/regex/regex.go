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
