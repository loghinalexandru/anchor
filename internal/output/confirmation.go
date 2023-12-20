package output

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/loghinalexandru/anchor/internal/output/bubbletea/style"
)

const maxRetries = 3

func Confirmation(prompt string, in io.Reader, out io.Writer, renderer style.RenderFunc) bool {
	reader := bufio.NewReader(in)
	retries := 0

	for retries < maxRetries {
		_, err := fmt.Fprint(out, renderer(fmt.Sprintf("%s [y/n]: ", prompt)))
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
			_, _ = fmt.Fprintln(out, renderer("Aborting..."))
			return false
		}

		retries++
	}

	_, _ = fmt.Fprintln(out, renderer("Exceeded retry count"))
	return false
}
