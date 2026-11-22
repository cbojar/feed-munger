package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/beevik/etree"
	"github.com/cbojar/feed-munger/configuration"
)

func main() {
	config, err := configuration.Read(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read command: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Output: %s, Inputs: %s\n", config.Output(), config.Input())

	if input, err := readFile(config.Input()); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	} else if document, err := parseFile(input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %s\n", err)
		os.Exit(1)
	} else if document, err := processFile(config, document); err != nil {
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

func parseFile(data []byte) (*etree.Document, error) {
	document := etree.NewDocument()
	if err := document.ReadFromBytes(data); err != nil {
		return nil, err
	} else {
		return document, nil
	}
}

func processFile(config *configuration.Configuration, document *etree.Document) (*etree.Document, error) {
	processed := document.Copy()

	if config.Alternate() != "" {
		if link := processed.FindElement("/feed/link[@rel='alternate']"); link == nil {
			link = processed.CreateElement("link")
			link.CreateAttr("rel", "alternate")
			link.CreateAttr("type", "text/html")
			link.CreateAttr("href", config.Alternate())
			processed.SelectElement("feed").InsertChildAt(0, link)
		} else if href := link.SelectAttr("href"); href == nil {
			link.CreateAttr("href", config.Alternate())
		} else {
			href.Value = config.Alternate()
		}
	}

	if config.Name() != "" {
		if title := processed.FindElement("/feed/title"); title == nil {
			title = processed.CreateElement("title")
			title.SetText(config.Name())
			processed.SelectElement("feed").InsertChildAt(0, title)
		} else {
			title.SetText(config.Name())
		}
	}

	if config.ID() != "" {
		if id := processed.FindElement("/feed/id"); id == nil {
			id = processed.CreateElement("id")
			id.SetText(config.ID())
			processed.SelectElement("feed").InsertChildAt(0, id)
		} else {
			id.SetText(config.ID())
		}
	}

	return processed, nil
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
