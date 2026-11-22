package main

import (
	"fmt"
	"io"
	"os"

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

func unparseFile(document *etree.Document) ([]byte, error) {
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
