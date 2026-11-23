package main

import (
	"fmt"
	"os"

	c "github.com/cbojar/feed-munger/configuration"
	"github.com/cbojar/feed-munger/feeds/in"
	"github.com/cbojar/feed-munger/feeds/out"
	"github.com/cbojar/feed-munger/inputs"
)

func main() {
	options, err := c.Read(os.Args[1:])
	check(err, "Failed to read command")

	fmt.Fprintf(os.Stderr, "%s\n", options)

	inputs, closer, err := inputs.OpenAll(options.Inputs())
	check(err, "Error reading files")
	defer closer()

	documents, err := in.ParseAll(inputs)
	check(err, "Error parsing files")

	document, err := processFiles(options, documents)
	check(err, "Error processing files")

	output, err := document.Bytes()
	check(err, "Error serializing file")

	err = writeFile(options.Output(), output)
	check(err, "Error writing file")
}

func check(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", message, err)
		os.Exit(1)
	}
}

func processFiles(options *c.Options, documents *in.Feeds) (*out.Feed, error) {
	processed := out.NewFeed(options, documents.First())
	processed.AddAllNamespaces(documents.Namespaces())
	processed.AddAllEntries(documents.Entries())

	return processed, nil
}

func writeFile(path string, data []byte) error {
	if path == "-" {
		_, err := os.Stdout.Write(data)
		return err
	} else {
		return os.WriteFile(path, data, 0644)
	}
}
