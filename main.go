package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/beevik/etree"
	c "github.com/cbojar/feed-munger/configuration"
)

func main() {
	config, err := c.Read(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read command: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "%s\n", config)

	if inputs, err := readFiles(config.Inputs()); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	} else if documents, err := parseFiles(inputs); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %s\n", err)
		os.Exit(1)
	} else if document, err := processFiles(config, documents); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing file: %s\n", err)
		os.Exit(1)
	} else if output, err := unparseFile(document); err != nil {
		fmt.Fprintf(os.Stderr, "Error serializing file: %s\n", err)
		os.Exit(1)
	} else if err = writeFile(config.Output(), output); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		os.Exit(1)
	}
}

func readFiles(paths []string) ([][]byte, error) {
	files := make([][]byte, 0, len(paths))

	for _, path := range paths {
		if file, err := readFile(path); err != nil {
			return nil, err
		} else {
			files = append(files, file)
		}
	}

	return files, nil
}

func readFile(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	} else if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		if response, err := http.Get(path); err != nil {
			return nil, err
		} else {
			defer response.Body.Close()
			return io.ReadAll(response.Body)
		}
	} else {
		return os.ReadFile(path)
	}
}

func parseFiles(datas [][]byte) ([]*etree.Document, error) {
	documents := make([]*etree.Document, 0, len(datas))

	for _, data := range datas {
		if document, err := parseFile(data); err != nil {
			return nil, err
		} else {
			documents = append(documents, document)
		}
	}

	return documents, nil
}

func parseFile(data []byte) (*etree.Document, error) {
	document := etree.NewDocument()
	if err := document.ReadFromBytes(data); err != nil {
		return nil, err
	} else {
		return document, nil
	}
}

func processFiles(config *c.Configuration, documents []*etree.Document) (*etree.Document, error) {
	processed := createFeed(config, documents[0])
	root := processed.Root()

	for _, document := range documents {
		for _, attr := range document.Root().Attr {
			if strings.HasPrefix(attr.FullKey(), "xmlns") {
				root.CreateAttr(attr.FullKey(), attr.Value)
			}
		}
	}

	entries := []*etree.Element{}
	for _, document := range documents {
		entries = append(entries, extractEntries(document)...)
	}

	slices.SortFunc(entries, func(entry1 *etree.Element, entry2 *etree.Element) int {
		date1 := entry1.SelectElement("updated").Text()
		date2 := entry2.SelectElement("updated").Text()

		return strings.Compare(date2, date1)
	})

	for _, entry := range entries {
		root.AddChild(entry)
	}

	return processed, nil
}

func createFeed(config *c.Configuration, firstDocument *etree.Document) *etree.Document {
	document := etree.NewDocument()
	document.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	root := document.CreateElement("feed")
	root.CreateAttr("xmlns", "http://www.w3.org/2005/Atom")

	id := root.CreateElement("id")
	if config.ID() != "" {
		id.SetText(config.ID())
	} else if firstDocumentID := firstDocument.FindElement("/feed/id"); firstDocumentID != nil {
		id.SetText(firstDocumentID.Text())
	}

	title := root.CreateElement("title")
	if config.Name() != "" {
		title.SetText(config.Name())
	} else if firstDocumentTitle := firstDocument.FindElement("/feed/title"); firstDocumentTitle != nil {
		title.SetText(firstDocumentTitle.Text())
	}

	alternate := root.CreateElement("link")
	alternate.CreateAttr("rel", "alternate")
	alternate.CreateAttr("type", "text/html")
	if config.Alternate() != "" {
		alternate.CreateAttr("href", config.Alternate())
	} else if firstDocumentAlternate := firstDocument.FindElement("/feed/link[@rel='alternate']"); firstDocumentAlternate != nil {
		alternate.CreateAttr("href", firstDocumentAlternate.SelectAttrValue("href", ""))
	}

	updated := root.CreateElement("updated")
	updated.SetText(time.Now().Format(time.RFC3339))

	return document
}

func extractEntries(document *etree.Document) []*etree.Element {
	feedTitle := document.Root().SelectElement("title").Text()
	entries := document.Root().SelectElements("entry")
	extracted := make([]*etree.Element, 0, len(entries))

	for _, entry := range entries {
		extractedEntry := entry.Copy()
		source := extractedEntry.CreateElement("source")
		sourceTitle := source.CreateElement("title")
		sourceTitle.SetText(feedTitle)
		extracted = append(extracted, extractedEntry)
	}

	return extracted
}

func unparseFile(document *etree.Document) ([]byte, error) {
	document.IndentTabs()
	return document.WriteToBytes()
}

func writeFile(path string, data []byte) error {
	if path == "-" {
		_, err := os.Stdout.Write(data)
		return err
	} else {
		return os.WriteFile(path, data, 0644)
	}
}
