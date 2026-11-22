package configuration

import (
	"flag"
	"fmt"
	"strings"
)

type Configuration struct {
	inputs []string
	output string
}

func Read(args []string) (*Configuration, error) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	output := flags.String("output", "-", "The destination to write the output feed to, or - for stdout")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	inputs := flags.Args()
	if len(inputs) == 0 {
		return nil, fmt.Errorf("No input files given")
	}

	return &Configuration{
		inputs: inputs,
		output: *output,
	}, nil
}

func (configuration *Configuration) Input() string {
	return configuration.inputs[0]
}

func (configuration *Configuration) Output() string {
	return configuration.output
}

func (configuration *Configuration) String() string {
	return fmt.Sprintf(
		"Configuration{output: %s, inputs: %s}",
		configuration.output,
		strings.Join(configuration.inputs, ", "),
	)
}
