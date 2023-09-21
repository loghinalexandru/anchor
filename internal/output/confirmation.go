package output

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Confirmation(s string, in io.Reader) bool {
	reader := bufio.NewReader(in)
	retryMax := 3

	for retryMax > 0 {
		fmt.Printf("%s [y/n]: ", s)

		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			fmt.Println("Aborting...")
			return false
		}

		retryMax--
	}

	fmt.Println("Exceeded retry count")
	return false
}
