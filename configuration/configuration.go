package configuration

import (
	"flag"
	"fmt"
	"strings"
)

type Configuration struct {
	inputs    []string
	output    string
	id        string
	name      string
	alternate string
}

func Read(args []string) (*Configuration, error) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	output := flags.String("output", "-", "The destination to write the output feed to, or - for stdout")
	id := flags.String("id", "", "The ID for the processed feed, or use the current ID of the feed if blank")
	name := flags.String("name", "", "The name for the processed feed, or use the current name of the feed if blank")
	alternate := flags.String("alternate", "", "The URL for the alternate representation of this feed, or use the current if left blank")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	inputs := flags.Args()
	if len(inputs) == 0 {
		return nil, fmt.Errorf("No input files given")
	}

	return &Configuration{
		inputs:    inputs,
		output:    *output,
		id:        *id,
		name:      *name,
		alternate: *alternate,
	}, nil
}

func (configuration *Configuration) Inputs() []string {
	return configuration.inputs
}

func (configuration *Configuration) Output() string {
	return configuration.output
}

func (configuration *Configuration) ID() string {
	return configuration.id
}

func (configuration *Configuration) Name() string {
	return configuration.name
}

func (configuration *Configuration) Alternate() string {
	return configuration.alternate
}

func (configuration *Configuration) String() string {
	return fmt.Sprintf(
		"Configuration{output: %s, id: %s, name: %s, alternate: %s, inputs: %s}",
		configuration.output,
		configuration.id,
		configuration.name,
		configuration.alternate,
		strings.Join(configuration.inputs, ", "),
	)
}
