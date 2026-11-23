package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/beevik/etree"
	c "github.com/cbojar/feed-munger/configuration"
	in "github.com/cbojar/feed-munger/feeds/in"
	"github.com/cbojar/feed-munger/inputs"
)

func main() {
	config, err := c.Read(os.Args[1:])
	check(err, "Failed to read command")

	fmt.Fprintf(os.Stderr, "%s\n", config)

	inputs, closer, err := inputs.OpenAll(config.Inputs())
	check(err, "Error reading files")
	defer closer()

	documents, err := in.ParseAll(inputs)
	check(err, "Error parsing files")

	document, err := processFiles(config, documents)
	check(err, "Error processing files")

	output, err := unparseFile(document)
	check(err, "Error serializing file")

	err = writeFile(config.Output(), output)
	check(err, "Error writing file")
}

func check(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", message, err)
		os.Exit(1)
	}
}

func processFiles(config *c.Configuration, documents *in.Feeds) (*etree.Document, error) {
	processed := createFeed(config, documents.First())
	root := processed.Root()

	for _, namespace := range documents.Namespaces() {
		root.CreateAttr(namespace.Key(), namespace.Value())
	}

	entries := []*etree.Element{}
	for _, document := range documents.All() {
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

func createFeed(config *c.Configuration, firstDocument *in.Feed) *etree.Document {
	document := etree.NewDocument()
	document.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	root := document.CreateElement("feed")
	root.CreateAttr("xmlns", "http://www.w3.org/2005/Atom")

	id := root.CreateElement("id")
	if config.ID() != "" {
		id.SetText(config.ID())
	} else if firstDocumentID, found := firstDocument.ID(); found {
		id.SetText(firstDocumentID)
	}

	title := root.CreateElement("title")
	if config.Name() != "" {
		title.SetText(config.Name())
	} else if firstDocumentTitle, found := firstDocument.Title(); found {
		title.SetText(firstDocumentTitle)
	}

	alternate := root.CreateElement("link")
	alternate.CreateAttr("rel", "alternate")
	alternate.CreateAttr("type", "text/html")
	if config.Alternate() != "" {
		alternate.CreateAttr("href", config.Alternate())
	} else if firstDocumentAlternate, found := firstDocument.Alternate(); found {
		alternate.CreateAttr("href", firstDocumentAlternate)
	}

	updated := root.CreateElement("updated")
	updated.SetText(time.Now().Format(time.RFC3339))

	return document
}

func extractEntries(document *in.Feed) []*etree.Element {
	feedTitle, titleFound := document.Title()

	entries := document.Entries()
	extracted := make([]*etree.Element, 0, len(entries))

	for _, entry := range entries {
		extractedEntry := entry.Copy()

		if titleFound {
			source := extractedEntry.CreateElement("source")
			sourceTitle := source.CreateElement("title")
			sourceTitle.SetText(feedTitle)
		}

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
