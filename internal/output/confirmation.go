package output

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Confirmation(s string, in io.Reader, out io.Writer) bool {
	reader := bufio.NewReader(in)
	retryMax := 3

	for retryMax > 0 {
		_, err := fmt.Fprintf(out, `%s [y/n]: `, s)
		if err != nil {
			return false
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		switch strings.ToLower(strings.TrimSpace(response)) {
		case "y", "yes":
			return true
		case "n", "no":
			_, _ = fmt.Fprintf(out, "Aborting...")
			return false
		}

		retryMax--
	}

	_, _ = fmt.Fprintf(out, "Exceeded retry count")
	return false
}
