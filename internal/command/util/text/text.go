package text

import (
	"fmt"
	"regexp"

	"github.com/loghinalexandru/anchor/internal/config"
)

func FindLines(content []byte, pattern string) [][]byte {
	regex := regexp.MustCompile(fmt.Sprintf(config.RegexpLine, regexp.QuoteMeta(pattern)))
	return regex.FindAll(content, -1)
}

func DeleteLines(content []byte, pattern string) []byte {
	regexPattern := regexp.MustCompile(fmt.Sprintf(config.RegexpLine, regexp.QuoteMeta(pattern)))
	return regexPattern.ReplaceAll(content, []byte(""))
}
