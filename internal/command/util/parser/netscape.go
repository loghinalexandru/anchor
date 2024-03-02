package parser

import (
	"errors"
	"os"
	"regexp"

	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/model"
	"github.com/virtualtam/netscape-go/v2"
)

var toolbarRegexp = regexp.MustCompile("(?i)bookmark|bar")

// TraverseNode creates files with appropriate labels based on the folder structure from the node.
// If there are invalid label names, they will be formatted according to label.Format function.
// Duplicate bookmarks are ignored by default and only the first occurrence is added in the file with
// the same label structure.
func TraverseNode(rootDir string, labels []string, node netscape.Folder) error {
	toolbar := toolbarRegexp.MatchString(node.Name)

	if len(node.Bookmarks) > 0 && !toolbar {
		labels = append(labels, node.Name)
	}

	err := createFile(rootDir, node.Bookmarks, labels)
	if err != nil {
		return err
	}

	for _, n := range node.Subfolders {
		err = TraverseNode(rootDir, labels, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func createFile(rootDir string, bookmarks []netscape.Bookmark, labels []string) error {
	if len(bookmarks) == 0 {
		return nil
	}

	file, err := label.Open(rootDir, label.Format(labels), os.O_APPEND|os.O_CREATE|os.O_RDWR)
	if err != nil {
		return err
	}

	for _, b := range bookmarks {
		entry, err := model.NewBookmark(b.URL, model.WithTitle(b.Title))
		if err != nil {
			return file.Close()
		}

		err = entry.Write(file)
		if err != nil && !errors.Is(err, model.ErrDuplicateBookmark) {
			return file.Close()
		}
	}

	return file.Close()
}
