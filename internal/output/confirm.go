package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/loghinalexandru/anchor/internal/output/bubbletea/style"
)

// Confirm is a wrapper function for Confirmer that uses os.Stdin for input,
// os.Stdout for output and a no operation Confirmer.Renderer.
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

// Confirm shows the user via out parameter a prompt in order to confirm the action. The input is
// read via the in parameter and keeps retrying until "yes/no" or "y/n" is found.
//
// Blocks until correct input is given or the number of retries exceeded MaxRetries.
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

	_, _ = fmt.Fprintln(out, c.Renderer("Exceeded retry count. Aborting..."))
	return false
}
