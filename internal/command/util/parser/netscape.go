package parser

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/model"
	"github.com/virtualtam/netscape-go/v2"
)

var toolbarRegexp = regexp.MustCompile("(?i)bookmark|bar")

func Traversal(rootDir string, labels []string, node netscape.Folder) error {
	bar := toolbarRegexp.MatchString(node.Name)

	if len(node.Bookmarks) > 0 && !bar {
		labels = append(labels, node.Name)
	}

	file, err := label.Open(rootDir, label.Format(labels), os.O_APPEND|os.O_CREATE|os.O_RDWR)
	if err != nil {
		return err
	}

	for _, b := range node.Bookmarks {
		entry, err := model.NewBookmark(b.URL, model.WithTitle(b.Title))
		if err != nil {
			file.Close()
			return err
		}

		err = entry.Write(file)
		if err != nil && !errors.Is(err, model.ErrDuplicateBookmark) {
			file.Close()
			return err
		}

		if errors.Is(err, model.ErrDuplicateBookmark) {
			fmt.Println(err)
		}
	}

	err = file.Close()
	if err != nil {
		return err
	}

	for _, n := range node.Subfolders {
		err = Traversal(rootDir, labels, n)
		if err != nil {
			return err
		}
	}

	return nil
}
