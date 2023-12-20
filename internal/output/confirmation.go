package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/loghinalexandru/anchor/internal/output/bubbletea/style"
)

func Confirm(prompt string) bool {
	return Confirmer{
		MaxRetries: 3,
		Renderer:   style.Nop,
	}.Confirm(prompt, os.Stdin, os.Stdout)
}

type Confirmer struct {
	MaxRetries int
	Renderer   style.RenderFunc
}

func (c Confirmer) Confirm(prompt string, in io.Reader, out io.Writer) bool {
	reader := bufio.NewReader(in)
	retries := 0

	for retries < c.MaxRetries {
		_, err := fmt.Fprint(out, c.Renderer(fmt.Sprintf("%s [y/n]: ", prompt)))
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
			_, _ = fmt.Fprintln(out, c.Renderer("Aborting..."))
			return false
		}

		retries++
	}

	_, _ = fmt.Fprintln(out, c.Renderer("Exceeded retry count"))
	return false
}
