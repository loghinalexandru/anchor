package text

import (
	"fmt"
	"regexp"
)

const regexLine = `(?i).+%s.+ .+\n`

func FindLines(content []byte, pattern string) [][]byte {
	regex := regexp.MustCompile(fmt.Sprintf(regexLine, regexp.QuoteMeta(pattern)))
	return regex.FindAll(content, -1)
}

func DeleteLines(content []byte, pattern string) []byte {
	regexPattern := regexp.MustCompile(fmt.Sprintf(regexLine, regexp.QuoteMeta(pattern)))
	return regexPattern.ReplaceAll(content, []byte(""))
}
